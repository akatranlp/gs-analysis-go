run:
	go run ./cmd/gs-analysis/main.go

build:
	go build -o ./bin/gs-analysis cmd/gs-analysis/main.go

dev:
	docker compose up

generate:
	sqlc generate

.PHONY: db-migrate
db-status:
	goose -dir ./internal/database/migrations sqlite3 ./db/db.sqlite3 status

.PHONY: db-up
db-up:
	goose -dir ./internal/database/migrations sqlite3 ./db/db.sqlite3 up

.PHONY: db-down
db-down:
	goose -dir ./internal/database/migrations sqlite3 ./db/db.sqlite3 down

up:
	docker-compose -f docker-compose-prod.yaml build
