# Сборка бинарей с приложением
FROM golang:1.20 AS build
COPY . /src
WORKDIR /src
# Для корректной работы требуется https://docs.docker.com/build/buildkit/#getting-started
# Необходимо для ускорения повторной сборки за счет использования кэша go-пакетов.
# Либо убрать `--mount=type=cache,target=/go`
RUN --mount=type=cache,target=/go make build

# Промежуточный образ, на основе которого будет собран финальный
FROM alpine:latest AS bin-image
COPY Makefile /app/Makefile
WORKDIR /app
RUN apk add gcompat make

# Final image stages
FROM bin-image AS app-image
COPY --from=build /src/bin/app /app/bin/app
CMD ["/app/bin/app"]

FROM bin-image AS http-image
COPY --from=build /src/bin/http /app/bin/http
CMD ["/app/bin/http"]
