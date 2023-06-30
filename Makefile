NAME=display_parser

### Build
.PHONY: image
image:
	sudo DOCKER_BUILDKIT=1 docker build -t ${NAME}-app:latest --target=app-image -f Dockerfile .
	sudo DOCKER_BUILDKIT=1 docker build -t ${NAME}-http:latest --target=http-image -f Dockerfile .

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build: vendor test
	go build -o bin/app ./cmd/app
	go build -o bin/http ./cmd/http

### Run
.PHONY: run
run:
	bin/app \
		--http-timeout=30s \
		--http-delay-per-request=3s \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=localhost \
		--db-port=5432 \
		--db-name=display_parser \
		--pipeline-use-stored-pages-only=false \
		--pipeline-model-parser-count=1 \
		--pipeline-page-collector-count=1

.PHONY: run-http
run-http:
	bin/http \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=localhost \
		--db-port=5432 \
		--db-name=display_parser


.PHONY: run-page-cache
run-page-cache:
	bin/app \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=localhost \
		--db-port=5432 \
		--db-name=display_parser \
		--pipeline-use-stored-pages-only=true \
		--pipeline-model-parser-count=1 \
		--pipeline-page-collector-count=1

.PHONY: run-swagger-ui
run-swagger-ui:
	sudo docker run --rm -p 80:8080 -e SWAGGER_JSON=/openapi.yml -v $$(pwd)/docs/openapi.yml:/openapi.yml swaggerapi/swagger-ui
### Dev
.PHONY: test
test: vendor
	go clean -testcache
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run -v

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
	sql-migrate up

