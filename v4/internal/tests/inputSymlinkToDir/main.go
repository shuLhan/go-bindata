// Copyright 2020 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

//go:generate go run main.go

//+build ignore

package main

import (
	"log"
	"regexp"

	"github.com/shuLhan/go-bindata/v4"
)

func main() {
	cfg := &bindata.Config{
		Package:     "bindata",
		AssetPrefix: bindata.DefAssetPrefixName,
		Prefix:      regexp.MustCompile(".*/testdata/"),
		ModTime:     1586263518,
		Input: []bindata.InputConfig{
			bindata.CreateInputConfig("../../../testdata/symlinkParent/..."),
		},
	}

	err := bindata.Translate(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
