BIN ?= dist/db

.PHONY: seed
seed:
	go build -tags "fts5" -v -o $(BIN) .
	$(BIN) seed

generate:
	sqlc generate

test:
	go test ./...

test.fixtures:
	UPDATE_FIXTURES=true make test

lint:
	golangci-lint run