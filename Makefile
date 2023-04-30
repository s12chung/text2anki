export KHAIII_BIN_PATH ?= integrations/tokenizers/dist/khaiii
export KOMORAN_JAR_PATH ?= integrations/tokenizers/dist/komoran
export TOKENIZER ?= khaiii
export DICTIONARY ?= krdict
BIN ?= dist/text2anki
include Makefile_env.mk

setup:
	cd integrations/tokenizers; make build

build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

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

test: test.diff
	go test -tags "$(TAGS)" $(TEST)
test.diff: db.diff
test.fixtures:
	# generate top level fixtures first
	UPDATE_FIXTURES=true go test $(TEST) -run TestGen___ || true
	UPDATE_FIXTURES=true make test

lint:
	golangci-lint run
goimports:
	goimports -w .

ci.build: build
ci.diff: test.diff
ci.test:
	go test -tags "$(TAGS)" -v $(TEST)
ci.setup:
	mkdir -p $(CI_BIN)
	cd db; make ci.setup

db.seed:
	cd db; make seed
db.generate:
	cd db; make generate
db.diff:
	cd db; make diff