package main

import (
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
)

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

func TestAsset(t *testing.T) {
	tests := []struct {
		desc   string
		name   string
		expErr string
		exp    string
	}{{
		desc:   "With invalid asset",
		name:   "in/split/test.1",
		expErr: "open in/split/test.1: file does not exist",
	}, {
		desc: "With valid asset",
		name: "in/a/test.asset",
		exp: `// sample file
`,
	}, {
		desc: "With space on asset",
		name: "in/file name",
		exp: `// Content of "testdata/in/file name"
`,
	}}

	for _, test := range tests {
		t.Log(test.desc, ":", test.name)

		got, err := Asset(test.name)
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.exp, string(got), true)
	}
}
