# Multi-stage build для оптимизации размера образа
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Установка CA сертификатов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование собранного приложения
COPY --from=builder /app/main .

# Копирование веб-файлов
COPY --from=builder /app/web ./web

# Открытие порта
EXPOSE 8080

# Команда запуска
CMD ["./main"]
