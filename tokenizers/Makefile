BUILD_DIR := build

build: clean
	mkdir -p $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/komoran; cd komoran; make jar; cp build/export/*.jar ../$(BUILD_DIR)/komoran

clean:
	rm -rf $(BUILD_DIR)
