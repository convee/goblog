SHELL := /bin/bash

APP := goblog
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%SZ)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)
PACKAGE_NAME := $(APP)-$(VERSION)-$(GOOS)-$(GOARCH)
PACKAGE_ROOT := dist/$(PACKAGE_NAME)
ARCHIVE := dist/$(PACKAGE_NAME).tar.gz

.PHONY: all
# make, Run checks and build for the current platform
all: check build

.PHONY: build
# make build, Build the application binary for GOOS/GOARCH
build:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o "bin/$(APP)" ./

.PHONY: build-admin
# make build-admin, Build the administrator bootstrap command
build-admin:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -o "bin/$(APP)-admin" ./cmd/admin

.PHONY: dep
# make dep, Download and verify dependencies
dep:
	@go mod download
	@go mod verify

.PHONY: test
# make test, Run tests with the race detector
test:
	@go test -race ./...

.PHONY: vet
# make vet, Run Go static analysis
vet:
	@go vet ./...

.PHONY: check
# make check, Run formatting, vet, and tests
check:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	@go vet ./...
	@go test -race ./...

.PHONY: vuln
# make vuln, Scan reachable Go vulnerabilities
vuln:
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: package
# make package, Build a complete release archive and checksum
package: build build-admin
	@rm -rf "$(PACKAGE_ROOT)"
	@mkdir -p "$(PACKAGE_ROOT)/bin" "$(PACKAGE_ROOT)/conf"
	@cp "bin/$(APP)" "bin/$(APP)-admin" "$(PACKAGE_ROOT)/bin/"
	@cp -R templates static migrations "$(PACKAGE_ROOT)/"
	@cp conf/config.example.yml "$(PACKAGE_ROOT)/conf/config.yml"
	@cp blog.sql startup.sh README.md SECURITY.md LICENSE CHANGELOG.md THIRD_PARTY_NOTICES.md "$(PACKAGE_ROOT)/"
	@tar -C dist -czf "$(ARCHIVE)" "$(PACKAGE_NAME)"
	@if command -v sha256sum >/dev/null 2>&1; then cd dist && sha256sum "$(notdir $(ARCHIVE))" > "$(notdir $(ARCHIVE)).sha256"; else cd dist && shasum -a 256 "$(notdir $(ARCHIVE))" > "$(notdir $(ARCHIVE)).sha256"; fi
	@echo "created $(ARCHIVE)"

.PHONY: tar
# make tar, Alias for make package
tar: package

.PHONY: fmt
# make fmt, Format Go source files
fmt:
	@gofmt -s -w .

.PHONY: clean
# make clean, Remove generated local build artifacts
clean:
	@rm -f "bin/$(APP)" "bin/$(APP)-admin"
	@rm -rf dist
	@echo "clean finished"

.PHONY: help
# make help, Show available targets
help:
	@echo 'Usage: make [target]'
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
