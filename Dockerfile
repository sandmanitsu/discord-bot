# Используем официальный образ Go для сборки
FROM golang:1.22 as builder

# Устанавливаем зависимости для сборки
RUN apt-get update && apt-get install -y build-essential libopus-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка бота
RUN go build -o discord-bot .

# Минимальный runtime-образ для выполнения бота
FROM ubuntu:latest

# Установка opus и ffmpeg библиотек для выполнения
RUN apt-get update && apt-get install -y ffmpeg libopus-dev

# Установка рабочей директории
WORKDIR /app

# Копируем собранный бинарник из builder-образа
COPY --from=builder /app/discord-bot .

# Копируем конфигурационные файлы
COPY config ./config
COPY google ./google

# Запуск бота
CMD ["./discord-bot"]