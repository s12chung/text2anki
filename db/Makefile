BIN ?= dist/db
include Makefile_env.mk

build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

.PHONY: seed
seed: build
	$(BIN) seed

generate:
	sqlc generate

TEST ?= ./...
test:
	go test -tags "$(TAGS)" $(TEST)

test.fixtures:
	UPDATE_FIXTURES=true make test

lint:
	golangci-lint run