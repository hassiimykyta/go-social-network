COMPOSE ?= docker compose
DB_DSN=postgres://postgres:postgres@localhost:5432/app?sslmode=disable

.PHONY: up
up: 
	$(COMPOSE) up --build

.PHONY: restart
restart: 
	$(COMPOSE) restart

.PHONY: stop
stop: 
	$(COMPOSE) stop

.PHONY: down
down: 
	$(COMPOSE) down

.PHONY: clean
clean: 
	$(COMPOSE) down -v

migrate-up:
	goose -dir ./db/migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir ./db/migrations postgres "$(DB_DSN)" down

migrate-status:
	goose -dir ./db/migrations postgres "$(DB_DSN)" status

generate-queries:
	rm -rf internal/dbgen && sqlc generate