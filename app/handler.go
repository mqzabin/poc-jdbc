package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createEventHandler(conn *pgxpool.Pool) http.HandlerFunc {
	handleErr := func(w http.ResponseWriter, err error) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			handleErr(w, err)
			return
		}

		type eventUserID struct {
			UserID string `json:"user_id"`
		}

		var msg eventUserID
		if err := json.Unmarshal(data, &msg); err != nil {
			handleErr(w, err)
			return
		}

		const query = "insert into events(event, user_id) values ($1, $2)"

		if _, err := conn.Exec(r.Context(), query, data, msg.UserID); err != nil {
			handleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
