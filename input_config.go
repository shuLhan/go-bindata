// Copyright 2020 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"path/filepath"
	"strings"
)

// InputConfig defines options on an asset directory to be convert.
type InputConfig struct {
	// Path defines a directory containing asset files to be included
	// in the generated output.
	Path string

	// Recursive defines whether subdirectories of path
	// should be recursively included in the conversion.
	Recursive bool
}

func CreateInputConfig(path string) InputConfig {
	inConfig := newInputConfig(path)
	return *inConfig
}

//
// newInputConfig determines whether the given path has a recursive indicator
// ("/...") and returns a new path with the recursive indicator chopped off if
// it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
//
func newInputConfig(path string) *InputConfig {
	inConfig := &InputConfig{}

	if strings.HasSuffix(path, "/...") {
		inConfig.Path = filepath.Clean(path[:len(path)-4])
		inConfig.Recursive = true
	} else {
		inConfig.Path = filepath.Clean(path)
	}

	return inConfig
}
