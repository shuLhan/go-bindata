// Copyright 2020 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"testing"
)

func TestNewInputConfig(t *testing.T) {
	tests := []struct {
		desc string
		path string
		exp  *InputConfig
	}{{
		desc: `With suffix /...`,
		path: `./...`,
		exp: &InputConfig{
			Path:      `.`,
			Recursive: true,
		},
	}, {
		desc: `Without suffix /...`,
		path: `.`,
		exp: &InputConfig{
			Path: `.`,
		},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		got := newInputConfig(test.path)

		assert(t, test.exp, got, true)
	}
}
