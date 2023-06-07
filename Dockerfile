FROM golang:1.20 AS build

COPY . /src
WORKDIR /src
RUN go mod vendor && go build -o display_parser .

FROM alpine:latest
RUN apk add gcompat
WORKDIR /app
COPY --from=build /src/display_parser ./
CMD ["./display_parser"]
