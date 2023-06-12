include ../Makefile_env.mk

mkfile_path := $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

export KHAIII_BIN_PATH ?= $(mkfile_path)/integrations/tokenizers/dist/khaiii
BIN ?= dist/db


build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

.PHONY: seed
seed: build
	$(BIN) seed
search: build
	$(BIN) search

generate:
	sqlc generate
	make build; $(BIN) generate

diff:
	sqlc diff
	make build; $(BIN) diff

TEST ?= ./...
test:
	go test -tags "$(TAGS)" $(TEST)

test.fixtures:
	UPDATE_FIXTURES=true make test

SQLC_VERSION = 1.18.0
SQLC_URL := https://github.com/kyleconroy/sqlc/releases/download/v$(SQLC_VERSION)/sqlc_$(SQLC_VERSION)_linux_amd64.tar.gz
ci.setup:
	curl -sSL $(SQLC_URL) | tar zxvf - -C ../$(CI_BIN)

lint:
	golangci-lint run