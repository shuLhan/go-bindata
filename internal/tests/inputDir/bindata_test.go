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
	cases := []struct {
		desc   string
		name   string
		expErr string
		exp    string
	}{{
		desc:   "With invalid asset",
		name:   "in/split",
		expErr: "open in/split: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/split/",
		expErr: "open in/split/: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/split/test.1",
		expErr: "open in/split/test.1: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/split/test.2",
		expErr: "open in/split/test.2: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/a/",
		expErr: "open in/a/: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/a",
		expErr: "open in/a: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/a/test.asset",
		expErr: "open in/a/test.asset: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/b/test.asset",
		expErr: "open in/b/test.asset: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/c/test.asset",
		expErr: "open in/c/test.asset: file does not exist",
	}, {
		desc: "With space on asset",
		name: "in/file name",
		exp: `// Content of "testdata/in/file name"
`,
	}, {
		desc: "With valid asset",
		name: "in/test.asset",
		exp: `// sample file
`,
	}}

	for _, c := range cases {
		t.Log(c.desc, ":", c.name)

		got, err := Asset(c.name)
		if err != nil {
			assert(t, c.expErr, err.Error(), true)
			continue
		}

		assert(t, c.exp, string(got), true)
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
