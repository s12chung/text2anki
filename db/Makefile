include ../Makefile_env.mk

mkfile_path := $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

export KHAIII_BIN_PATH ?= $(mkfile_path)/integrations/tokenizers/dist/khaiii
BIN ?= dist/db


build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

create: build
	$(BIN) create
.PHONY: seed
seed: build
	$(BIN) seed
testdb: build
	$(BIN) testdb
search: build
	$(BIN) search

generate:
	sqlc generate

diff:
	sqlc diff

TEST ?= ./...
test: testdb
	go test -tags "$(TAGS)" $(TEST)

test.fixtures: testdb
	# generate top level fixtures first
	UPDATE_FIXTURES=true go test $(TEST) -run TestGen___ | $(FIXTURE_CLEAN_OUTPUT) || true
	@echo; echo
	UPDATE_FIXTURES=true make test | $(FIXTURE_CLEAN_OUTPUT)

SQLC_VERSION = 1.18.0
SQLC_URL := https://github.com/kyleconroy/sqlc/releases/download/v$(SQLC_VERSION)/sqlc_$(SQLC_VERSION)_linux_amd64.tar.gz
ci.setup:
	curl -sSL $(SQLC_URL) | tar zxvf - -C ../$(CI_BIN)

lint:
	golangci-lint run