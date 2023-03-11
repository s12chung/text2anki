JDK_PATH := /Library/Java/JavaVirtualMachines/jdk1.8.0_101.jdk/Contents/Home
BIN := dist/text2anki
export JAVA_HOME := $(JDK_PATH)/jre
export CGO_CFLAGS := -I$(JDK_PATH)/include -I$(JDK_PATH)/include/darwin

setup:
	cd tokenizers; make build

lint:
	golangci-lint run

goimports:
	goimports -w .

build:
	go build -o $(BIN) .

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

test:
	go test ./...

test.fixtures:
	UPDATE_FIXTURES=true make test
