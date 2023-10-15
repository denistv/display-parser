NAME=display_parser

### Build
.PHONY: image
image:
	sudo DOCKER_BUILDKIT=1 docker build -t ${NAME}-app:latest --target=app-image -f Dockerfile .
	sudo DOCKER_BUILDKIT=1 docker build -t ${NAME}-http:latest --target=http-image -f Dockerfile .

.PHONY: image-dev-tools
image-dev-tools:
	sudo DOCKER_BUILDKIT=1 docker build -f Dockerfile-dev-tools -t display_parser-dev-tools .

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build:
	go build -race -o bin/app ./cmd/app
	go build -race -o bin/http ./cmd/http

### Run
.PHONY: run
run:
	sudo docker-compose run app \
        /usr/local/bin/app \
		--http-timeout=30s \
		--http-delay-per-request=3s \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=postgres \
		--db-port=5432 \
		--db-name=display_parser \
		--pipeline-use-stored-pages-only=false \
		--pipeline-model-parser-count=1 \
		--pipeline-page-collector-count=1

.PHONY: run-cached
run-cached:
	sudo docker-compose run app \
        /usr/local/bin/app \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=postgres \
		--db-port=5432 \
		--db-name=display_parser \
		--pipeline-pages-cache=true \
		--pipeline-model-parser-count=8 \
		--pipeline-page-collector-count=1

.PHONY: run-swagger-ui
run-swagger-ui:
	sudo docker run --rm -p 80:8080 -e SWAGGER_JSON=/openapi.yml -v $$(pwd)/docs/openapi.yml:/openapi.yml swaggerapi/swagger-ui
### Dev
.PHONY: test
test:
	go test -race -v ./...

.PHONY: lint
lint:
	mkdir -p ~/.cache/golangci-lint
	sudo docker run \
		--rm \
		-v ${HOME}/.cache/golangci-lint:/home/$$(id -u -n)/.cache/golangci-lint \
		-v $$(pwd):/src \
		-w /src \
		-t \
		--user $$(id -u):$$(id -g) \
		-e HOME=/home/$$(id -u -n) \
		golangci/golangci-lint:v1.53.3 \
		mkdir -p /home/$$(id -u -n) && golangci-lint run -v

.PHONY: install-dev-tools
install-dev-tools:
	go install github.com/vektra/mockery/v2@v2.20.0
	go install github.com/rubenv/sql-migrate/...@v1.4.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.1
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: mock
mock:
	mockery --dir internal

.PHONY: run-docker
run-docker:
	sudo docker-compose up

.PHONY: migrate
migrate:
	sudo docker-compose run sql-migrate

.PHONY: all
all: mock lint test image
