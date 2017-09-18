##
## This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain
## Dedication license. Its contents can be found at:
## http://creativecommons.org/publicdomain/zero/1.0
##

.PHONY: all build
.PHONY: lint lint-errors lint-all
.PHONY: test coverbrowse test-cmd
.PHONY: clean distclean

CMD_DIR           :=./go-bindata
TESTDATA_DIR      :=./testdata/
TESTDATA_IN_DIR   :=./testdata/in
TESTDATA_OUT_DIR  :=./testdata/out

SRC               :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .GoFiles }} {{ $$dir }}/{{.}} {{end}}' ./...)
TEST              :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .TestGoFiles }} {{ $$dir }}/{{.}} {{end}}' ./...)
TEST_COVER_ALL    :=cover.out

TARGET_CMD        :=$(shell go list -f '{{ .Target }}' $(CMD_DIR))
CMD_SRC           :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .GoFiles }} {{ $$dir }}/{{.}} {{end}}' $(CMD_DIR))
CMD_TEST          :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .TestGoFiles }} {{ $$dir }}/{{.}} {{end}}' $(CMD_DIR))
CMD_COVER_OUT     :=$(CMD_DIR)/cover.out

TARGET_LIB        :=$(shell go list -f '{{ .Target }}' ./)
LIB_SRC           :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .GoFiles }} {{ $$dir }}/{{.}} {{end}}' ./)
LIB_TEST          :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .TestGoFiles }} {{ $$dir }}/{{.}} {{end}}' ./)
LIB_COVER_OUT     :=lib.cover.out

TEST_COVER_OUT    :=$(TEST_COVER_ALL) $(LIB_COVER_OUT) $(CMD_COVER_OUT)
TEST_COVER_HTML   :=cover.html

TEST_OUT          := \
	$(TESTDATA_OUT_DIR)/compress-memcopy.go \
	$(TESTDATA_OUT_DIR)/compress-nomemcopy.go \
	$(TESTDATA_OUT_DIR)/debug.go \
	$(TESTDATA_OUT_DIR)/nocompress-memcopy.go \
	$(TESTDATA_OUT_DIR)/nocompress-nomemcopy.go

VENDOR_DIR        :=$(PWD)/vendor
VENDOR_BIN        :=$(VENDOR_DIR)/bin
LINTER_CMD        :=$(VENDOR_BIN)/gometalinter
LINTER_DEF_OPTS   :=--vendor --concurrency=1 --disable=gotype --deadline=240s
LINTER            :=GOBIN=$(VENDOR_BIN) $(VENDOR_BIN)/gometalinter $(LINTER_DEF_OPTS)

##
## MAIN TARGET
##

all: build test-cmd

##
## CLEAN
##

clean:
	rm -f $(TEST_COVER_OUT) $(TEST_COVER_HTML)

distclean: clean
	rm -rf $(TARGET_CMD) $(TARGET_LIB) $(VENDOR_DIR)

##
## LINT
##

$(VENDOR_DIR): vendor.deps
	@echo ">>> Installing vendor dependencies ..."
	@./.scripts/deps.sh $<

$(LINTER_CMD): $(VENDOR_DIR)

lint: $(LINTER_CMD)
	@echo ">>> Linting ..."
	@$(LINTER) --fast ./...

lint-errors: $(LINTER_CMD)
	@echo ""
	@echo ">>> Lint errors only ..."
	@$(LINTER) --fast --errors ./...

lint-all: $(LINTER_CMD)
	@echo ">>> Run all linters ..."
	@$(LINTER) --exclude="testdata/*" ./...

##
## TEST
##

$(LIB_COVER_OUT): $(LIB_SRC) $(LIB_TEST)
	@echo ""
	@echo ">>> Testing library ..."
	@go test -v -coverprofile=$@ ./

$(CMD_COVER_OUT): $(CMD_SRC) $(CMD_TEST)
	@echo ""
	@echo ">>> Testing cmd ..."
	@go test -v -coverprofile=$@ $(CMD_DIR)

$(TEST_COVER_ALL): $(LIB_COVER_OUT) $(CMD_COVER_OUT)
	@echo ""
	@echo ">>> Generate single coverage '$@' ..."
	@cat $^ | sed '/mode: set/d' | sed '1s/^/mode: set\n/' > $@

$(TEST_COVER_HTML): $(TEST_COVER_ALL)
	@echo ">>> Generate HTML coverage '$@' ..."
	@go tool cover -html=$< -o $@

test: lint-errors $(TEST_COVER_HTML)

coverbrowse: test
	@xdg-open $(TEST_COVER_HTML)

##
## TEST POST BUILD
##

$(TESTDATA_OUT_DIR)/compress-memcopy.go: $(TESTDATA_IN_DIR)/*
	$(TARGET_CMD) -o $@ -prefix="/.*/testdata/" $(TESTDATA_IN_DIR)/...
	@$(LINTER) --fast --errors $@

$(TESTDATA_OUT_DIR)/compress-nomemcopy.go: $(TESTDATA_IN_DIR)/*
	$(TARGET_CMD) -o $@ -prefix="/.*/testdata/" -nomemcopy $(TESTDATA_IN_DIR)/...
	@$(LINTER) --fast --errors $@

$(TESTDATA_OUT_DIR)/debug.go: $(TESTDATA_IN_DIR)/*
	$(TARGET_CMD) -o $@ -prefix="/.*/testdata/" -debug $(TESTDATA_IN_DIR)/...
	@$(LINTER) --fast --errors $@

$(TESTDATA_OUT_DIR)/nocompress-memcopy.go: $(TESTDATA_IN_DIR)/*
	$(TARGET_CMD) -o $@ -prefix="/.*/testdata/" -nocompress $(TESTDATA_IN_DIR)/...
	@$(LINTER) --fast --errors $@

$(TESTDATA_OUT_DIR)/nocompress-nomemcopy.go: $(TESTDATA_IN_DIR)/*
	$(TARGET_CMD) -o $@ -prefix="/.*/testdata/" -nocompress -nomemcopy $(TESTDATA_IN_DIR)/...
	@$(LINTER) --fast --errors $@

$(TEST_OUT): $(LINTER_CMD) $(TARGET_CMD)

test-cmd: $(TEST_OUT)

##
## BUILD
##

$(TARGET_LIB): $(LIB_SRC)
	go install ./

$(TARGET_CMD): $(CMD_SRC) $(TARGET_LIB)
	go install $(CMD_DIR)

build: test $(TARGET_CMD)
