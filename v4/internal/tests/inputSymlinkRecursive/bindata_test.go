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
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent",
		expErr: "open symlinkRecursiveParent: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent/symlinkTarget",
		expErr: "open symlinkRecursiveParent/symlinkTarget: file does not exist",
	}, {
		desc: "With valid asset",
		name: "symlinkRecursiveParent/file1",
		exp:  "// symlinkRecursiveParent/file1\n",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent/symlinkTarget/file1",
		expErr: "open symlinkRecursiveParent/symlinkTarget/file1: file does not exist",
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
