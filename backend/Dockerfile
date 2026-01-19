# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .

# Копируем конфигурации (если есть)
COPY --from=builder /app/config ./config

# Создаем пользователя для безопасности
RUN adduser -D -g '' appuser
USER appuser

# Запускаем приложение
EXPOSE 8080
CMD ["./main"]