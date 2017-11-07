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
TESTDATA_DIR      :=./testdata
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

POST_TEST_FILES   := \
	$(TESTDATA_DIR)/bindata_test.go \
	$(TESTDATA_DIR)/split_test.go

TEST_OUT          := \
	$(TESTDATA_OUT_DIR)/opt/no-output/bindata.go \
	$(TESTDATA_OUT_DIR)/compress/memcopy/bindata.go \
	$(TESTDATA_OUT_DIR)/compress/nomemcopy/bindata.go \
	$(TESTDATA_OUT_DIR)/debug/bindata.go \
	$(TESTDATA_OUT_DIR)/nocompress/memcopy/bindata.go \
	$(TESTDATA_OUT_DIR)/nocompress/nomemcopy/bindata.go \
	$(TESTDATA_OUT_DIR)/split/bindata.go

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
	rm -rf $(TEST_COVER_OUT) $(TEST_COVER_HTML) $(TESTDATA_OUT_DIR)

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
	@go test -v ./ && \
		go test -coverprofile=$@ ./ &>/dev/null && \
		go tool cover -func=$@


$(CMD_COVER_OUT): $(CMD_SRC) $(CMD_TEST)
	@echo ""
	@echo ">>> Testing cmd ..."
	@go test -v $(CMD_DIR) && \
		go test -coverprofile=$@ $(CMD_DIR) &>/dev/null && \
		go tool cover -func=$@

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

$(TESTDATA_OUT_DIR)/opt/no-output/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing without '-o' flag"
	@mkdir -p $(TESTDATA_OUT_DIR)/opt/no-output && \
		cd $(TESTDATA_OUT_DIR)/opt/no-output && \
		$(TARGET_CMD) -pkg main -prefix="/.*/testdata/" \
			-ignore="split/" ../../../../$(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go $(TESTDATA_OUT_DIR)/opt/no-output/
	@$(LINTER) $(TESTDATA_OUT_DIR)/opt/no-output/
	go test -v $(TESTDATA_OUT_DIR)/opt/no-output/

$(TESTDATA_OUT_DIR)/compress/memcopy/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing default option (compress, memcopy)"
	@$(TARGET_CMD) -o $@ -pkg main -prefix="/.*/testdata/" \
		-ignore="split/" $(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go $(TESTDATA_OUT_DIR)/compress/memcopy/
	@$(LINTER) $(TESTDATA_OUT_DIR)/compress/memcopy/
	go test -v $(TESTDATA_OUT_DIR)/compress/memcopy/

$(TESTDATA_OUT_DIR)/compress/nomemcopy/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing with '-nomemcopy'"
	@$(TARGET_CMD) -o $@ -pkg main -prefix="/.*/testdata/" \
		-ignore="split/" -nomemcopy $(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go \
		$(TESTDATA_OUT_DIR)/compress/nomemcopy/
	@$(LINTER) $(TESTDATA_OUT_DIR)/compress/nomemcopy/
	go test -v $(TESTDATA_OUT_DIR)/compress/nomemcopy/

$(TESTDATA_OUT_DIR)/debug/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing opt 'debug'"
	@$(TARGET_CMD) -o $@ -pkg main -prefix="/.*/testdata/" \
		-ignore="split/" -debug $(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go $(TESTDATA_OUT_DIR)/debug/
	@$(LINTER) $(TESTDATA_OUT_DIR)/debug/
	go test -v $(TESTDATA_OUT_DIR)/debug/

$(TESTDATA_OUT_DIR)/nocompress/memcopy/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing opt '-nocompress'"
	@$(TARGET_CMD) -o $@ -pkg main -prefix="/.*/testdata/" \
		-ignore="split/" -nocompress $(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go $(TESTDATA_OUT_DIR)/nocompress/memcopy/
	@$(LINTER) $(TESTDATA_OUT_DIR)/nocompress/memcopy/
	go test -v $(TESTDATA_OUT_DIR)/nocompress/memcopy/

$(TESTDATA_OUT_DIR)/nocompress/nomemcopy/bindata.go: $(TESTDATA_IN_DIR)/*
	@echo ">>> Testing opt '-nocompress -nomemcopy'"
	@$(TARGET_CMD) -o $@ -pkg main -prefix="/.*/testdata/" \
		-ignore="split/" -nocompress -nomemcopy $(TESTDATA_IN_DIR)/...
	@cp $(TESTDATA_DIR)/bindata_test.go $(TESTDATA_OUT_DIR)/nocompress/nomemcopy/
	@$(LINTER) $(TESTDATA_OUT_DIR)/nocompress/nomemcopy/
	go test -v $(TESTDATA_OUT_DIR)/nocompress/nomemcopy/

$(TESTDATA_OUT_DIR)/split/bindata.go: $(TESTDATA_IN_DIR)/split/*
	@echo ">>> Testing opt '-split'"
	@$(TARGET_CMD) -o $(TESTDATA_OUT_DIR)/split/ -pkg main \
		-prefix="/.*/testdata/" -split $(TESTDATA_IN_DIR)/split/...
	@cp $(TESTDATA_DIR)/split_test.go $(TESTDATA_OUT_DIR)/split/
	@$(LINTER) $(TESTDATA_OUT_DIR)/split/
	go test -v $(TESTDATA_OUT_DIR)/split/

$(TEST_OUT): $(LINTER_CMD) $(TARGET_CMD) $(POST_TEST_FILES)

test-cmd: $(TEST_OUT)

##
## BUILD
##

$(TARGET_LIB): $(LIB_SRC)
	go install ./

$(TARGET_CMD): $(CMD_SRC) $(TARGET_LIB)
	go install $(CMD_DIR)

build: test $(TARGET_CMD)
