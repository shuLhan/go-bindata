// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"io/ioutil"
	"testing"
)

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
		desc:   "With invalid asset",
		name:   "testdata/in/test.asset",
		expErr: "open testdata/in/test.asset: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/",
		expErr: "open in/: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in",
		expErr: "open in: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/test",
		expErr: "open in/test: file does not exist",
	}, {
		desc: "With valid asset",
		name: "in/test.asset",
		exp:  "// sample file\n",
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

func TestGeneratedContent(t *testing.T) {
	expFile := "bindata.exp"
	gotFile := "bindata.go"

	// Compare the generate bindata.go with expected.
	exp, err := ioutil.ReadFile(expFile)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(gotFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(exp, got) {
		t.Fatalf("%s not match with %s", expFile, gotFile)
	}
}
