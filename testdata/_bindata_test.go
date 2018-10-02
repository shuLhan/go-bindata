// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
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
		name:   "in/split/",
		expErr: "open in/split/: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/split",
		expErr: "open in/split: file does not exist",
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
		name:   "in/split/test.3",
		expErr: "open in/split/test.3: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/split/test.4",
		expErr: "open in/split/test.4: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/a/",
		expErr: "open in/a/: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "in/a",
		expErr: "open in/a: file does not exist",
	}, {
		desc: "With valid asset",
		name: "in/a/test.asset",
		exp:  "// sample file\n",
	}, {
		desc: "With valid asset",
		name: "in/b/test.asset",
		exp:  "// sample file\n",
	}, {
		desc: "With valid asset",
		name: "in/c/test.asset",
		exp:  "// sample file\n",
	}, {
		desc: "With space on asset",
		name: "in/file name",
		exp:  "// Content of \"testdata/in/file name\"\n",
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
