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
OUTPUT_FILE ?= tmp/$(shell date +"%Y-%m-%d_%H-%M-%S")

run: build
	mkdir -p tmp
	test -e $(INPUT_FILE) || echo "This is a sample file. Put your own text at: $(INPUT_FILE)\n\n김밥입니다." > $(INPUT_FILE)
	$(BIN) $(INPUT_FILE) $(OUTPUT_FILE)

test:
	go test ./...

test.fixtures:
	UPDATE_FIXTURES=true make test
