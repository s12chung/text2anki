include Makefile_env.mk

mkfile_path := $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

export KHAIII_BIN_PATH ?= $(mkfile_path)/integrations/tokenizers/dist/khaiii
export KOMORAN_JAR_PATH ?= $(mkfile_path)/integrations/tokenizers/dist/komoran
export TOKENIZER ?= khaiii
export DICTIONARY ?= krdict
BIN ?= dist/text2anki

run: build
	mkdir -p tmp
	$(BIN)

open:
	(make run & (cd ui; sleep 1; npm run open)) | tee tmp/text2anki.log

setup:
	cd integrations/tokenizers; make build

build:
	go build -tags "$(TAGS)" -v -o $(BIN) .

TEST ?= ./...

test: test.diff
	go test -tags "$(TAGS)" $(TEST)
test.nocache:
	go test -count=1 -tags "$(TAGS)" $(TEST)
test.diff: db.diff

test.fixtures:
	# generate top level fixtures first
	UPDATE_FIXTURES=true go test $(TEST) -run TestGen___ | $(FIXTURE_CLEAN_OUTPUT) || true
	@echo; echo
	UPDATE_FIXTURES=true make test | $(FIXTURE_CLEAN_OUTPUT)
test.slow:
	go test -v -count=1 -json -tags "$(TAGS)" $(TEST) \
	| jq -r 'select(.Action == "pass" and .Test != null) | (.Package | split("/") | last ) + "," + .Test + "," + (.Elapsed | tostring)' \
	| sort -k3 -n -t, \
	| tail -n 25

lint:
	golangci-lint run $(TEST)
lint.fix:
	goimports -w .

ci.build: build
ci.diff: test.diff
ci.test:
	go test -v -tags "$(TAGS)" $(TEST)
ci.setup:
	mkdir -p $(CI_BIN)
	cd db; make ci.setup

db.seed:
	cd db; make seed
db.generate:
	cd db; make generate
db.diff:
	cd db; make diff