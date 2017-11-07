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
		desc       string
		assetName  string
		expContent string
		expErr     string
	}{{
		desc:       `With asset "in/split/test.1"`,
		assetName:  `in/split/test.1`,
		expContent: "// sample file 1\n",
	}, {
		desc:       `With asset "in/split/test.2"`,
		assetName:  `in/split/test.2`,
		expContent: "// sample file 2\n",
	}, {
		desc:       `With asset "in/split/test.3"`,
		assetName:  `in/split/test.3`,
		expContent: "// sample file 3\n",
	}, {
		desc:       `With asset "in/split/test.4"`,
		assetName:  `in/split/test.4`,
		expContent: "// sample file 4\n",
	}, {
		desc:      `With non existing asset "in/split/test.5"`,
		assetName: `in/split/test.5`,
		expErr:    "open in/split/test.5: file does not exist",
	}}

	for _, test := range tests {
		t.Log(test.desc)

		got, err := Asset(test.assetName)
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.expContent, string(got), true)
	}
}
