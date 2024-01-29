POC_ENV?="poc.env"

.PHONY: setup/cleanup
setup/cleanup:
	@docker compose --env-file=$(POC_ENV) --profile=setup stop
	@docker compose --env-file=$(POC_ENV) --profile=setup rm -f

.PHONY: run/cleanup
run/cleanup:
	@docker compose --env-file=$(POC_ENV) --profile=test stop
	@docker compose --env-file=$(POC_ENV) --profile=test rm -f

.PHONY: cleanup
cleanup: run/cleanup setup/cleanup

.PHONY: setup
setup: run/cleanup setup/cleanup
	@cd app && go build -o ./bin *.go
	@docker compose --env-file=$(POC_ENV) --profile=setup up --build --detach
	@rm app/bin

.PHONY: run
run:
	@docker compose --env-file=$(POC_ENV) --profile=test up --build
	@echo "Access PGAdmin (email: admin@pg.com, password: admin) at http://localhost:80"
	@echo "Access KafDrop at http://localhost:9000"

.PHONY: fullrun
fullrun: setup run