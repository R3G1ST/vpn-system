FROM golang:1.21-alpine

WORKDIR /app

# Копируем только необходимые файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Сборка
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
