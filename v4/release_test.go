// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import "testing"

// nolint: gochecknoglobals
var sanitizeTests = []struct {
	in  string
	out string
}{
	{`hello`, "`hello`"},
	{"hello\nworld", "`hello\nworld`"},
	{"`ello", "(\"`\" + `ello`)"},
	{"`a`e`i`o`u`", "(((\"`\" + `a`) + (\"`\" + (`e` + \"`\"))) + ((`i` + (\"`\" + `o`)) + (\"`\" + (`u` + \"`\"))))"},
	{"\xEF\xBB\xBF`s away!", "(\"\\xEF\\xBB\\xBF\" + (\"`\" + `s away!`))"},
}

func TestSanitize(t *testing.T) {
	for _, tt := range sanitizeTests {
		out := sanitize([]byte(tt.in))
		if string(out) != tt.out {
			t.Errorf("sanitize(%q):\nhave %q\nwant %q", tt.in, out, tt.out)
		}
	}
}
