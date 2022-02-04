.PHONY: build

TAG=$(shell git describe --abbrev=0 --tags)
DATE=$(shell go run ./scripts/date.go)

build:
		@go mod tidy && \
		go build -ldflags "-X main.version=$(TAG) -X main.buildDate=$(DATE)" -o tran

install: tran
		@mv tran /usr/local/bin
