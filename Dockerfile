# Сборка бинарей с приложением
FROM golang:1.20 AS build
COPY . /src
WORKDIR /src
# Для корректной работы требуется https://docs.docker.com/build/buildkit/#getting-started
# Необходимо для ускорения повторной сборки за счет использования кэша go-пакетов.
# Благодаря этому компиляция и пересборка итоговых образов занимает считанные секунды.
# Либо убрать `--mount=type=cache,target=/go`
RUN --mount=type=cache,target=/go make vendor
RUN --mount=type=cache,target=/root/.cache/go-build make build

# Промежуточный образ, на основе которого будет собран финальный
FROM alpine:3.18.2 AS bin-image
COPY Makefile /app/Makefile
WORKDIR /app
RUN --mount=type=cache,target=/var/cache/apk apk add gcompat make

# Final image stages
FROM bin-image AS app-image
COPY --from=build /src/bin/app /usr/local/bin/app
CMD ["/usr/local/bin/app"]

FROM bin-image AS http-image
COPY --from=build /src/bin/http /usr/local/bin/http
CMD ["/usr/local/bin/http"]
