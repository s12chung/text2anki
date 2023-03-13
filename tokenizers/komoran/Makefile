KOMORAN_VERSION := 3.3.7
CLASS_PATH := lib/KOMORAN-$(KOMORAN_VERSION).jar lib/json-20210307.jar
TARGET_DIR := target
DIST_DIR := dist

joinwith = $(subst $(eval) ,$1,$2)

build: clean
	mkdir -p $(TARGET_DIR)
	javac -d $(TARGET_DIR) -classpath $(call joinwith,:,$(CLASS_PATH)): $(shell find ./src/main/java -name *.java) 

JAR_MANIFEST := Manifest.txt
JAR_NAME := tokenizer-komoran.jar

jar: build
	cd $(TARGET_DIR); \
		mkdir -p lib; cp ../lib/*.jar lib; \
		echo "Main-Class: text2anki.tokenizer.komoran.Server" > $(JAR_MANIFEST); \
		echo "Class-Path: $(call joinwith, ,$(CLASS_PATH))" >> $(JAR_MANIFEST); \
		jar cfmv $(JAR_NAME) $(JAR_MANIFEST) text2anki; \
		mkdir -p ../$(DIST_DIR)/lib; cp $(JAR_NAME) ../$(DIST_DIR); cp lib/*.jar ../$(DIST_DIR)/lib;

server: jar
	cd $(TARGET_DIR); java -jar $(JAR_NAME);

clean:
	rm -rf $(TARGET_DIR) $(DIST_DIR)

TEST_NAME := TestTokenizer
TEST_ARGS := -classpath "$(DIST_DIR)/$(JAR_NAME):$(call joinwith,:,$(CLASS_PATH)):"

test: jar
	javac -d $(TARGET_DIR) $(TEST_ARGS) $(shell find ./src/test/java -name *.java) 
	cd $(TARGET_DIR); java $(TEST_ARGS) $(TEST_NAME);
