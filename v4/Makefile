## Copyright 2018 The go-bindata Authors. All rights reserved.
## Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
## Public Domain Dedication license that can be found in the LICENSE file.

.PHONY: all build test clean distclean

CMD_DIR       :=./cmd/go-bindata

SRC           :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .GoFiles }} {{ $$dir }}/{{.}} {{end}}' ./...)
SRC_TEST      :=$(shell go list -f '{{ $$dir := .Dir }}{{ range .TestGoFiles }} {{ $$dir }}/{{.}} {{end}}' ./...)

TEST_COVER_OUT    :=cover.out
TEST_COVER_HTML   :=cover.html

TEST_LIB := \
	internal/tests/inputDir/bindata.go \
	internal/tests/inputDirRecursive/bindata.go \
	internal/tests/inputDuplicateDir/bindata.go \
	internal/tests/inputDuplicateFile/bindata.go \
	internal/tests/inputFile/bindata.go \
	internal/tests/inputSymlinkRecursive/bindata.go \
	internal/tests/inputSymlinkToDir/bindata.go \
	internal/tests/inputSymlinkToFile/bindata.go \
	internal/tests/withDebug/bindata.go \
	internal/tests/withNoCompress/bindata.go \
	internal/tests/withNoCompressNoMemCopy/bindata.go \
	internal/tests/withNoMemCopy/bindata.go \
	internal/tests/withSplit/bindata.go \
	internal/tests/withoutOutputFlag/bindata.go

##
## MAIN TARGET
##

all: build

##
## CLEAN
##

clean:
	rm -rf $(TEST_COVER_OUT) $(TEST_COVER_HTML) $(TEST_LIB)

distclean: clean
	go clean -i ./...

##
## TEST
##

%/bindata.go: %/main.go %/bindata_test.go $(LIB_SRC)
	go generate $<

$(TEST_COVER_OUT): $(SRC) $(SRC_TEST) $(TEST_LIB)
	@echo ">>> Testing ..."
	go test -coverprofile=$@ ./...

$(TEST_COVER_HTML): $(TEST_COVER_OUT)
	@echo ">>> Generate HTML coverage '$@' ..."
	@go tool cover -html=$< -o $@

test: $(TEST_COVER_HTML)

##
## BUILD
##

build: test
	go build ./...
