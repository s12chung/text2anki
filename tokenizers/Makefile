DIST_DIR := dist

build: clean komoran

.PHONY: khaiii
khaiii:
	rm -rf $(DIST_DIR)/khaiii;
	mkdir -p $(DIST_DIR)/khaiii; cd khaiii; make dist; cp -r $(DIST_DIR)/* ../$(DIST_DIR)/khaiii	

.PHONY: komoran
komoran:
	rm -rf $(DIST_DIR)/komoran;
	mkdir -p $(DIST_DIR)/komoran; cd komoran; make jar; cp -r $(DIST_DIR)/* ../$(DIST_DIR)/komoran	

clean:
	rm -rf $(DIST_DIR)
