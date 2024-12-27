# vim: ts=4 sw=4
# Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License 
# Gorilla
#
SHELL := /bin/bash
.DEFAULT_GOAL := all

BASE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SRC_DIR := $(BASE_DIR)
BIN_DIR := $(BASE_DIR)bin/
CMD_DIR := $(SRC_DIR)cmd/
EXE_EXT := 

# collect all source files (excluding tests)
GO_SRCS=$(shell find $(SRC_DIR) -name '*.go' -not -name '*_test.go')
GO_CMDS=$(shell find $(CMD_DIR) -name '*.go' -not -name '*_test.go')
GO_BINS=$(patsubst $(CMD_DIR)%.go,$(BIN_DIR)%,$(GO_CMDS))

GO_TEST_OPTS ?=

.PHONY: all info bins clean bins-init test unittest smoketest
info:
	@echo "BASE_DIR:   $(BASE_DIR)"
	@echo "SRC_DIR:    $(SRC_DIR)"
	@echo "BIN_DIR:    $(BIN_DIR)"
	@echo "CMD_DIR:    $(CMD_DIR)" 
	@echo "GO_SRCS:    $(GO_SRCS)"
	@echo "GO_CMDS:    $(GO_CMDS)"
	@echo "GO_BINS:    $(GO_BINS)"
	@echo


all: bins test
bins: bins-init $(GO_BINS) 

clean: bins-clean
bins-clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BASE_DIR)tmp

bins-init:
	mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && go mod tidy
	cd $(SRC_DIR) && go vet ./...

install:
	cp $(BIN_DIR)* $(HOME)/bin

$(BIN_DIR)%$(EXE_EXT): $(GO_SRCS) $(SRC_DIR)go.mod $(SRC_DIR)go.sum $(LIBRARY_JSONS)
	cd $(SRC_DIR) && go build -o $(BIN_DIR) $(patsubst $(BIN_DIR)%$(EXE_EXT),$(CMD_DIR)%.go,$@)

test: unittest 

unittest:
	mkdir -p $(BASE_DIR)tmp/coverage
	cd $(SRC_DIR) && go test $(GO_TEST_OPTS) -coverprofile=$(BASE_DIR)tmp/coverage/coverage.out.tmp $(SRC_DIR)...
	cat $(BASE_DIR)tmp/coverage/coverage.out.tmp \
		| grep -v "github.com/tommika/gorilla/assert/" \
		| grep -v "github.com/tommika/gorilla/must/" \
		> $(BASE_DIR)tmp/coverage/coverage.out
	cd $(SRC_DIR) && go tool cover -html=$(BASE_DIR)tmp/coverage/coverage.out -o=$(BASE_DIR)tmp/coverage/coverage.html
	cd $(SRC_DIR) && go tool cover -func=$(BASE_DIR)tmp/coverage/coverage.out -o=$(BASE_DIR)tmp/coverage/coverage.txt

bench:
	cd $(SRC_DIR) && go test -test.run=^X -benchmem -bench=. $(GO_TEST_OPTS) $(SRC_DIR)...

push:
	git add .
	sleep 1
	git commit -m "updates" . || true
	git push
