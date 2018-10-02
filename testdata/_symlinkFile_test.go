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
		exp    string
		expErr string
	}{{
		desc: "With valid asset",
		name: "symlinkFile/file1",
		exp:  "// symlink file 1\n",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkFile/file5",
		expErr: "open symlinkFile/file5: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkSrc/file1",
		expErr: "open symlinkSrc/file1: file does not exist",
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
