BIN := dist/khaiiserver

lint:
	golangci-lint run

build:
	go build -v -o $(BIN) .


PORT := 9999
run: build
	$(BIN) --port $(PORT)

POST_URI := http://localhost:$(PORT)/tokenize
STR := 안녕! 반가워!
curl:
	 curl -X POST $(POST_URI) -H 'Content-Type: application/json' -d '{"string": "$(STR)"}'

test:
	go test -v ./...

test.fixtures:
	UPDATE_FIXTURES=true make test