COMPOSE ?= docker compose

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