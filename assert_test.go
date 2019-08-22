// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
)

//nolint: gochecknoglobals
var (
	_traces = make([]byte, 1024)
	lerr    = log.New(os.Stderr, "", 0)
)

func printStack() {
	var lines, start, end int

	runtime.Stack(_traces, false)

	for x, b := range _traces {
		if b != '\n' {
			continue
		}

		lines++
		if lines == 5 {
			start = x + 1
		} else if lines == 7 {
			end = x + 1
			break
		}
	}

	lerr.Println("!!! ERR " + string(_traces[start:end]))
}

func assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) == equal {
		return
	}

	printStack()

	t.Fatalf("\n"+
		">>> Expecting '%+v'\n"+
		"          got '%+v'\n", exp, got)
	os.Exit(1)
}
