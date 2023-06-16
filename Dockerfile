# Собираем приложение
FROM golang:1.20 AS build

COPY . /src
WORKDIR /src

RUN make build

# Кладем артефакт с бинарником в финальный легковесный образ
FROM alpine:latest

RUN apk add gcompat make
WORKDIR /app
COPY --from=build /src/display_parser ./
COPY Makefile /app/Makefile

CMD ["make", "run"]
