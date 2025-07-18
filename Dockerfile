# 1. Базовый образ с Go
FROM golang:1.24-alpine AS builder

# 2. Устанавливаем рабочую папку
WORKDIR /app

# 3. Копируем go.mod/go.sum и скачиваем зависимости (кэшируем слои)
COPY go.mod go.sum ./
RUN go mod download

# 4. Копируем весь проект
COPY . .

# 5. Собираем бинарник (оптимизированная сборка)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/server/main.go

# 6. Финальный образ — минимальный
FROM alpine:latest
WORKDIR /root/

# 7. Копируем бинарник
COPY --from=builder /app/server .

# 8. Указываем порт, соответствующий вашему конфигу (например, 8080 и grpc-port)
EXPOSE 8080
EXPOSE 50051

# 9. Запуск сервера
CMD ["./server"]
