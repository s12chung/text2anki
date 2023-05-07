DIST_DIR := dist
BIN := $(DIST_DIR)/khaiiiserver

lint:
	golangci-lint run

build:
	go build -v -o $(BIN) .

PORT := 9999
run: build
	$(BIN) --port $(PORT)

POST_URI := http://localhost:$(PORT)/tokenize
STR ?= 안녕! 반가워!
curl:
	 curl -X POST $(POST_URI) -H 'Content-Type: application/json' -d '{"string": "$(STR)"}'

test:
	go test ./...

test.fixtures:
	UPDATE_FIXTURES=true make test

.PHONY: dist
dist: clean build
	cp -r lib $(DIST_DIR)

clean:
	rm -rf $(DIST_DIR)