### Build
.PHONY: image
image:
	sudo docker build -t display_parser:latest -f Dockerfile .

.PHONY: build
build:
	go build -o display_parser

### Run
.PHONY: run
run:
	./display_parser \
		--http-timeout=30s \
		--http-delay-per-request=500ms \
		--db-user=display_parser \
		--db-password=display_parser \
		--db-hostname=localhost \
		--db-port=5432 \
		--db-name=display_parser \
		--pipeline-use-stored-pages-only=false \
		--pipeline-model-parser-count=5 \
		--pipeline-page-collector-count=10

### Dev
.PHONY: test
test:
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
