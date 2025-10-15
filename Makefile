# Makefile для Chat Server

.PHONY: build run clean test deps

# Переменные
BINARY_NAME=chat-server
BUILD_DIR=build
MAIN_FILE=main.go

# Сборка приложения
build:
	@echo "Сборка приложения..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Приложение собрано: $(BUILD_DIR)/$(BINARY_NAME)"

# Запуск приложения
run:
	@echo "Запуск приложения..."
	@go run $(MAIN_FILE)

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	@go mod tidy
	@go mod download

# Очистка
clean:
	@echo "Очистка..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Тестирование
test:
	@echo "Запуск тестов..."
	@go test ./...

# Форматирование кода
fmt:
	@echo "Форматирование кода..."
	@go fmt ./...

# Проверка кода
vet:
	@echo "Проверка кода..."
	@go vet ./...

# Линтинг
lint:
	@echo "Линтинг..."
	@golangci-lint run

# Полная проверка
check: fmt vet test

# Запуск в режиме разработки с hot reload
dev:
	@echo "Запуск в режиме разработки..."
	@air

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build  - Собрать приложение"
	@echo "  run    - Запустить приложение"
	@echo "  deps   - Установить зависимости"
	@echo "  clean  - Очистить файлы сборки"
	@echo "  test   - Запустить тесты"
	@echo "  fmt    - Форматировать код"
	@echo "  vet    - Проверить код"
	@echo "  lint   - Запустить линтер"
	@echo "  check  - Полная проверка (fmt + vet + test)"
	@echo "  dev    - Запуск в режиме разработки"
	@echo "  help   - Показать эту справку"
