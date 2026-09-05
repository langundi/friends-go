include .env
MIGRATIONS_PATH = ./cmd/migrate

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-stop
docker-stop:
	docker compose stop

.PHONY: migration-create
migration-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migration-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migration-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migration-force
migration-force:
ifndef version
	$(error version is not set. Usage: make migration-force version=<N>)
endif
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(version)
