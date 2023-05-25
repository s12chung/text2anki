BIN ?= dist/db
include ../Makefile_env.mk

build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

.PHONY: seed
seed: build
	$(BIN) seed
search: build
	$(BIN) search

generate: build
	$(BIN) generate
	sqlc generate

diff: build
	sqlc diff
	$(BIN) diff

TEST ?= ./...
test:
	go test -tags "$(TAGS)" $(TEST)

test.fixtures:
	UPDATE_FIXTURES=true make test

SQLC_URL := https://github.com/kyleconroy/sqlc/releases/download/v1.17.2/sqlc_1.17.2_linux_amd64.tar.gz
ci.setup:
	curl -sSL $(SQLC_URL) | tar zxvf - -C ../$(CI_BIN)

lint:
	golangci-lint run