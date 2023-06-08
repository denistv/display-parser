.PHONY: image
image:
	sudo docker build -t display_parser .

.PHONY: run
run:
	sudo docker run display_parser

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	sudo docker run -t --rm -v $$(pwd):$$(pwd) -w $$(pwd) golangci/golangci-lint:latest golangci-lint run -v

.PHONY: install-dev-tools
install-dev-tools:
	go install github.com/vektra/mockery/v2@v2.20.0

.PHONY: mock
mock:
	mockery --dir internal