# script/Makefile

DB_URL ?= postgres://postgres:yourpassword@localhost:5432/edu_platform?sslmode=disable
SWAGGER_DIRS := ./cmd,./internal/user/transport/http,./internal/user/entity
MIGRATIONS_DIR := ./migrations

.PHONY: swagger run tidy migrate-up migrate-down migrate-create deploy help

swagger:
	swag init --dir $(SWAGGER_DIRS) --output ./docs

run:
	@echo "🚀 Запуск сервера разработки..."
	air

tidy:
	go mod tidy

migrate-up:
	@echo "🔼 Применение миграций..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@echo "🔽 Откат миграций..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-create:
	@test -n "$(name)" || (echo "❌ Укажи имя миграции: make migrate-create name=create_users"; exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

deploy:
	$(info 🚀 Деплой на $(SERVER_IP)...)
	@test -n "$(SERVER_IP)" || { echo "❌ Переменная SERVER_IP не задана"; exit 1; }
	@test -n "$(SERVER_PATH)" || { echo "❌ Переменная SERVER_PATH не задана"; exit 1; }
	@ssh root@$(SERVER_IP) "cd $(SERVER_PATH) && \
	 git pull && \
	 docker compose build backend && \
	 docker compose up -d --no-deps backend && \
	 echo '✅ Бэкенд успешно задеплоен!'"

help:
	$(info 🔧 Доступные команды:)
	$(info - make swagger        Сгенерировать Swagger документацию)
	$(info - make run            Запуск dev-сервера)
	$(info - make tidy           Очистить зависимости)
	$(info - make migrate-up     Применить миграции)
	$(info - make migrate-down   Откатить миграции)
	$(info - make migrate-create name=...  Создать миграцию)
	$(info - make deploy         Задеплоить на сервер)
	@exit 0
