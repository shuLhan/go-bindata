// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain
// Dedication license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"github.com/shuLhan/go-bindata"
)

var (
	_traces = make([]byte, 1024)
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
			start = x
		} else if lines == 7 {
			end = x + 1
			break
		}
	}

	os.Stderr.Write(_traces[start:end])
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

func TestParseArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	var (
		defConfig    = bindata.NewConfig()
		argInputPath = "."
		argPkg       = "pkgnametest"
		argOutPkg    = "assets"
		argOutFile   = argOutPkg + "/template.go"
	)

	tests := []struct {
		desc      string
		args      []string
		expErr    error
		expConfig *bindata.Config
	}{{
		desc: `Without input`,
		args: []string{
			"noop",
		},
		expErr: ErrNoInput,
	}, {
		desc: `With "-prefix prefix/*/to/be/removed ."`,
		args: []string{
			"noop",
			"-prefix", "prefix/*/to/be/removed",
			".",
		},
		expConfig: &bindata.Config{
			Output:  defConfig.Output,
			Package: "main",
			Prefix:  regexp.MustCompile("prefix/*/to/be/removed"),
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}, {
		desc: `With "-pkg ` + argPkg + `"`,
		args: []string{
			"noop",
			"-pkg", argPkg,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  defConfig.Output,
			Package: argPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}, {
		desc: `With "-o ` + argOutFile + `" (package name should be "` + argOutPkg + `")`,
		args: []string{
			"noop",
			"-o", argOutFile,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  argOutFile,
			Package: argOutPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}, {

		desc: `With "-pkg ` + argPkg + ` -o ` + argOutPkg + `" (package name should be "` + argPkg + `")`,
		args: []string{
			"noop",
			"-pkg", argPkg,
			"-o", argOutFile,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:  argOutFile,
			Package: argPkg,
			Input: []bindata.InputConfig{{
				Path: argInputPath,
			}},
			Ignore: defConfig.Ignore,
		},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		os.Args = test.args

		flag.CommandLine = flag.NewFlagSet(test.args[0],
			flag.ContinueOnError)

		initArgs()

		gotErr := parseArgs()

		if test.expErr != nil {
			assert(t, test.expErr, gotErr, true)
			continue
		}

		assert(t, test.expConfig, cfg, true)
	}
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		desc string
		path string
		exp  bindata.InputConfig
	}{{
		desc: `With suffix /...`,
		path: `./...`,
		exp: bindata.InputConfig{
			Path:      `.`,
			Recursive: true,
		},
	}, {
		desc: `Without suffix /...`,
		path: `.`,
		exp: bindata.InputConfig{
			Path: `.`,
		},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		got := parseInput(test.path)

		assert(t, test.exp, got, true)
	}
}
