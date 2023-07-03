TAGS ?= sqlite_math_functions
CI_BIN ?= .github/bin
PATH := $(PWD)/$(CI_BIN):$(PATH)

FIXTURE_CLEAN_OUTPUT = awk 'BEGIN{ORS=""; RS=""} {gsub(/\n[^\n]*    fixture\.go([^\n]*\n){2,6}[^\n]*turn off ENV var to run test/, " - UPDATE_FIXTURES=true")}1'
