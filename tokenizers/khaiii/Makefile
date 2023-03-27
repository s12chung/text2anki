BIN := dist/khaiiserver

lint:
	golangci-lint run

build:
	go build -v -o $(BIN) .

run: build
	$(BIN) --port 3000

test:
	go test -v ./...

test.fixtures:
	UPDATE_FIXTURES=true make test