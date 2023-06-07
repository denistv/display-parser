.PHONY: build-docker-image
build-docker-image:
	sudo docker build -t display_parser .

.PHONY: run
run:
	sudo docker run display_parser

.PHONY: test
test:
	sudo docker run -t --rm -v $$(pwd):$$(pwd) -w $$(pwd) golang:1.20 bash -c "go test ./..."

.PHONY: lint
lint:
	sudo docker run -t --rm -v $$(pwd):$$(pwd) -w $$(pwd) golangci/golangci-lint:latest golangci-lint run -v

.PHONY: all
all:
	make build-docker-image
	make lint
	make test
