JDK_PATH := /Library/Java/JavaVirtualMachines/jdk1.8.0_101.jdk/Contents/Home
export JAVA_HOME := $(JDK_PATH)/jre
export CGO_CFLAGS := -I$(JDK_PATH)/include -I$(JDK_PATH)/include/darwin


setup:
	cd tokenizers; make build

lint:
	golangci-lint run

goimports:
	goimports -w .

run:
	go run .

test:
	go test ./...

test.fixtures:
	UPDATE_FIXTURES=true make test || true # ensure true to run below
	cp pkg/dictionary/koreanbasic/testdata/search.xml pkg/anki/testdata/koreanbasic.xml
	false
