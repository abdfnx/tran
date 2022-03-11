.PHONY: build

TAG=$(shell git describe --abbrev=0 --tags)
DATE=$(shell go run ./scripts/date.go)

build:
		@go mod tidy -compat=1.17 && \
		go build -ldflags "-X main.version=$(TAG) -X main.buildDate=$(DATE)" -o tran

install: tran
		@mv tran /usr/local/bin

jbtc: # just build tran container without pushing it
		@docker build --file ./docker/vm/Dockerfile -t trancli/tran .

btc: # build tran container
		@docker push trancli/tran

btcwc: # build tran container with cache
		@docker pull trancli/tran:latest && \
		docker build -t trancli/tran --cache-from trancli/tran:latest . && \
		docker push trancli/tran

jbftc: # just build full tran container without pushing it
		@docker build --file ./docker/container/Dockerfile -t trancli/tran-full .

bftc: # build full tran container
		@docker push trancli/tran-full

bftcwc: # build full tran container with cache
		@docker pull trancli/tran-full:latest && \
		docker build -t trancli/tran-full --cache-from trancli/tran-full:latest . && \
		docker push trancli/tran-full

ght:
		@node ./scripts/gh-tran/gh-trn.js
