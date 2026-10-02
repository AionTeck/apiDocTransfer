COMPOSE := docker compose
SERVICE := application

.DEFAULT_GOAL := help

.PHONY: help build up down restart rebuild logs shell run stop once local

help: ## Показать список команд
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Собрать образ
	$(COMPOSE) build

up: ## Поднять контейнер в фоне
	$(COMPOSE) up -d

down: ## Остановить и удалить контейнер
	$(COMPOSE) down

restart: ## Перезапустить (подхватит изменения .env)
	$(COMPOSE) up -d --force-recreate

rebuild: ## Пересобрать образ и пересоздать контейнер
	$(COMPOSE) up -d --build --force-recreate

logs: ## Логи контейнера
	$(COMPOSE) logs -f $(SERVICE)

shell: ## Зайти в контейнер
	$(COMPOSE) exec $(SERVICE) sh

run: ## Запустить app в работающем контейнере: make run ARGS="..."
	$(COMPOSE) exec -d $(SERVICE) sh -c 'exec ./app "$$@" >/proc/1/fd/1 2>/proc/1/fd/2' app $(ARGS)

stop: ## Корректно остановить запущенный app (SIGINT, как Ctrl+C)
	$(COMPOSE) exec $(SERVICE) pkill -INT -x app || echo "app не запущен"
