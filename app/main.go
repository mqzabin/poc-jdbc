package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	DatabaseName     string `conf:"env:DATABASE_NAME,default:postgres"`
	DatabaseUser     string `conf:"env:DATABASE_USER,default:postgres"`
	DatabasePassword string `conf:"env:DATABASE_PASSWORD,default:postgres"`
	DatabaseHost     string `conf:"env:DATABASE_HOST,default:postgres"`
	DatabasePort     string `conf:"env:DATABASE_PORT,default:5432"`

	HTTPServerAddress         string        `conf:"env:HTTP_SERVER_ADDRESS,default:0.0.0.0:3000"`
	HTTPServerShutdownTimeout time.Duration `conf:"env:HTTP_SERVER_SHUTDOWN_TIMEOUT,default:10s"`
}

func main() {
	if err := Main(); err != nil {
		log.Println(err.Error())
	}
}

func Main() error {
	var config Env

	if _, err := conf.Parse("", &config); err != nil {
		return fmt.Errorf("parsing environment variables: %w", err)
	}

	ctx := context.Background()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.DatabaseUser, config.DatabasePassword,
		config.DatabaseHost, config.DatabasePort, config.DatabaseName,
	)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/create_event", createEventHandler(conn))

	server := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: mux,
	}

	shutdownCh := make(chan os.Signal)
	signal.Notify(shutdownCh, os.Interrupt, os.Kill)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("listening for http request:", err)
		}
	}()

	<-shutdownCh

	ctx, cancel := context.WithTimeout(context.Background(), config.HTTPServerShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down http server: %w", err)
	}

	wg.Done()

	return nil
}
