DIST_DIR := dist

build: clean
	mkdir -p $(DIST_DIR)
	mkdir -p $(DIST_DIR)/komoran; cd komoran; make jar; cp $(DIST_DIR)/*.jar ../$(DIST_DIR)/komoran

clean:
	rm -rf $(DIST_DIR)
