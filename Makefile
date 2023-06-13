.PHONY: image
image:
	sudo docker build -t display_parser:latest -f Dockerfile .

.PHONY: run
run:
	sudo docker run display_parser

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
