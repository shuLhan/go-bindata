// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

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
			Output:      defConfig.Output,
			Package:     "main",
			AssetPrefix: bindata.DefAssetPrefixName,
			Prefix:      regexp.MustCompile("prefix/*/to/be/removed"),
			Input: []bindata.InputConfig{
				bindata.CreateInputConfig(argInputPath),
			},
			Ignore:  defConfig.Ignore,
			Include: defConfig.Include,
		},
	}, {
		desc: `With "-pkg ` + argPkg + `"`,
		args: []string{
			"noop",
			"-pkg", argPkg,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:      defConfig.Output,
			Package:     argPkg,
			AssetPrefix: bindata.DefAssetPrefixName,
			Input: []bindata.InputConfig{
				bindata.CreateInputConfig(argInputPath),
			},
			Ignore:  defConfig.Ignore,
			Include: defConfig.Include,
		},
	}, {
		desc: `With "-o ` + argOutFile + `" (package name should be "` + argOutPkg + `")`,
		args: []string{
			"noop",
			"-o", argOutFile,
			argInputPath,
		},
		expConfig: &bindata.Config{
			Output:      argOutFile,
			Package:     argOutPkg,
			AssetPrefix: bindata.DefAssetPrefixName,
			Input: []bindata.InputConfig{
				bindata.CreateInputConfig(argInputPath),
			},
			Ignore:  defConfig.Ignore,
			Include: defConfig.Include,
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
			Output:      argOutFile,
			Package:     argPkg,
			AssetPrefix: bindata.DefAssetPrefixName,
			Input: []bindata.InputConfig{
				bindata.CreateInputConfig(argInputPath),
			},
			Ignore:  defConfig.Ignore,
			Include: defConfig.Include,
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
