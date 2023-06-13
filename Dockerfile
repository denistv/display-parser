# Тянем вендоров
FROM golang:1.20 AS vendor
COPY . /src
WORKDIR /src

RUN go mod vendor

# Собираем приложение
FROM golang:1.20 AS build

COPY --from=vendor /src /src
WORKDIR /src

RUN go build -o display_parser .

# Кладем артефакт с бинарником в финальный легковесный образ
FROM alpine:latest
RUN apk add gcompat
WORKDIR /app
COPY --from=build /src/display_parser ./
CMD ["./display_parser"]
