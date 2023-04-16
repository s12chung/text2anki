export KHAIII_BIN_PATH ?= integrations/tokenizers/dist/khaiii
export KOMORAN_JAR_PATH ?= integrations/tokenizers/dist/komoran
export TOKENIZER ?= khaiii
BIN ?= dist/text2anki

setup:
	cd integrations/tokenizers; make build

lint:
	golangci-lint run

goimports:
	goimports -w .

build:
	go build -v -o $(BIN) .

INPUT_FILE ?= tmp/in.txt
DEFAULT_INPUT_FILE := "이것은 샘플 파일입니다. $(INPUT_FILE)에 자신의 텍스트를 입력합니다.\n\nThis is a sample file. Put your own text at: $(INPUT_FILE)."
OUTPUT_DIR ?= tmp/$(shell date +"%Y-%m-%d_%H-%M-%S")

tmp:
	mkdir -p tmp
	test -e $(INPUT_FILE) || echo $(DEFAULT_INPUT_FILE) > $(INPUT_FILE)

run: build tmp
	$(BIN) --clean-speaker $(INPUT_FILE) $(OUTPUT_DIR)

subconv: tmp
	go run ./cmd/subconv $(INPUT_FILE) tmp/subconv.txt

syncfiltered:
	go run ./cmd/syncfiltered "$(SYNC_FILTERED_DIR)"


TEST ?= ./...
test:
	go test $(TEST)

ci.test:
	go test -v $(TEST)

test.fixtures:
	UPDATE_FIXTURES=true make test

db.seed:
	cd db; make seed
db.generate:
	cd db; make generate