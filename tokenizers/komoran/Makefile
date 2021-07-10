SRC := src/text2anki/tokenizer/komoran/*

KOMORAN_VERSION := 3.3.7
CLASS_PATH := libs/KOMORAN-$(KOMORAN_VERSION).jar libs/json-20210307.jar
BUILD_DIR := build

joinwith = $(subst $(eval) ,$1,$2)

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	javac -d $(BUILD_DIR) -classpath $(call joinwith,:,$(CLASS_PATH)): $(SRC).java

JAR_MANIFEST := Manifest.txt
JAR_NAME := tokenizer-komoran.jar

jar: build
	cd $(BUILD_DIR); \
		mkdir -p libs; cp ../libs/*.jar libs; \
		echo "Class-Path: $(call joinwith, ,$(CLASS_PATH))" > $(JAR_MANIFEST); \
		jar cfmv $(JAR_NAME) $(JAR_MANIFEST) text2anki; \
		mkdir -p export; cp $(JAR_NAME) export; cp ../libs/*.jar export;

clean:
	rm -rf $(BUILD_DIR)

TEST_NAME := TestTokenizer
TEST_ARGS := -classpath "build/$(JAR_NAME):$(call joinwith,:,$(CLASS_PATH)):"

test: clean jar
	javac -d $(BUILD_DIR) $(TEST_ARGS) $(TEST_NAME).java
	cd $(BUILD_DIR); java $(TEST_ARGS) $(TEST_NAME);
