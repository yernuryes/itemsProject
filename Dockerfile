# Используем базовый образ golang для сборки приложения
FROM golang:1.24-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходные файлы приложения
COPY . .

# Собираем приложение
WORKDIR /app/main
RUN go build -o main .

# Финальный образ на основе Alpine для уменьшения размера
FROM alpine:latest

# Копируем исполняемый файл приложения
COPY --from=build /app/main/main /app/main

# Копируем миграции базы данных
COPY --from=build /app/db/migrations /app/db/migrations

# Устанавливаем команду по умолчанию для запуска приложения
CMD ["/app/main"]

