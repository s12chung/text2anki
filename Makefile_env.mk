TAGS ?= sqlite_math_functions
CI_BIN ?= .github/bin
PATH := $(PWD)/$(CI_BIN):$(PATH)

_FIXTURE_CLEAN = 'BEGIN{ORS=""; RS=""} {gsub(/\n[^\n]*    fixture\.go([^\n]*\n){2,6}[^\n]*turn off ENV var to run test/, " - UPDATE_FIXTURES=true")}1'
_TESTDB_CLEAN = 'BEGIN{ORS=""; RS=""} {gsub(/\n[^\n]*    testdb_test\.go([^\n]*\n){2,6}[^\n]*turn off ENV var to run test/, " - UPDATE_FIXTURES=true")}1'

FIXTURE_CLEAN_OUTPUT = awk $(_FIXTURE_CLEAN) | awk $(_TESTDB_CLEAN)
