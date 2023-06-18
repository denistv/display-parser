# Сборка бинарей с приложением
FROM golang:1.20 AS build
COPY . /src
WORKDIR /src
RUN make build

# Промежуточный образ, на основе которого будет собран финальный
FROM alpine:latest AS bin-image
COPY Makefile /app/Makefile
WORKDIR /app
RUN apk add gcompat make

# Final image stages
FROM bin-image AS app-image
COPY --from=build /src/bin/app /app/app
LABEL bin=app
CMD ["make", "run"]

FROM bin-image AS http-image
COPY --from=build /src/bin/http /app/http
LABEL bin=http
CMD ["make", "run-http"]
