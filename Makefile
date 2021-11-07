SHELL := /bin/bash
BASEDIR = $(shell pwd)

COMMIT = $(shell git rev-parse --short HEAD)
TIME = $(shell TZ=Asia/Shanghai date +%Y%m%d%H)

.PHONY: build
# make build, Build the binary file
build: dep
	GOOS=linux GOARCH=amd64 go build -o "goblog" -v -ldflags "-X main.Commit=$COMMIT"

.PHONY: dep
# make dep Get the dependencies
dep:
	@go mod tidy

.PHONY: tar
# pack file
tar:
	@tar zcvf  goblog-"${TIME}".tar.gz tpl/ static/ goblog conf/ *.sh

.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .

.PHONY: clean
# make clean
clean:
	@-rm -vrf goblog*
	@go mod tidy
	@echo "clean finished"

# show help
help:
	@echo 'Usage: make [target]'
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all

