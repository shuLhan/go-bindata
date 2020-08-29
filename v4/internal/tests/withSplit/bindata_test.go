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
		exp    string
		expErr string
	}{{
		desc: "With valid asset",
		name: "in/split/test.1",
		exp:  "// sample file 1\n",
	}, {
		desc: "With valid asset",
		name: "in/split/test.2",
		exp:  "// sample file 2\n",
	}, {
		desc:   "With invalid asset",
		name:   "in/split/test.3",
		expErr: "open in/split/test.3: file does not exist",
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

func TestGeneratedAssets(t *testing.T) {
	exps := []string{
		"bindata",
		"bindataInSplitTest1",
		"bindataInSplitTest2",
	}

	for _, name := range exps {
		expFile := name + ".exp"
		gotFile := name + ".go"
		exp, err := ioutil.ReadFile(expFile)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ioutil.ReadFile(gotFile)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(exp, got) {
			t.Fatalf("%q not match with %q", expFile, gotFile)
		}
	}
}
