// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"os"
	"path/filepath"
	"unicode"
)

//
// asset holds information about a single asset to be processed.
//
type asset struct {
	// path contains full file path.
	path string

	// name contains key used in TOC -- name by which asset is referenced.
	name string

	// Function name for the procedure returning the asset contents.
	funcName string

	// fi field contains the file information (to minimize calling os.Stat
	// on the same file while processing).
	fi os.FileInfo
}

func normalize(in string) (out string) {
	up := true
	for _, r := range in {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if up {
				out += string(unicode.ToUpper(r))
				up = false
			} else {
				out += string(r)
			}
			continue
		}
		if r == '/' || r == '.' {
			up = true
		}
	}
	return out
}

//
// newAsset will create, initialize, and return new asset based on file
// path or real path if its symlink.
//
func newAsset(cfg *Config, path, name, realPath string, fi os.FileInfo) (ast *asset) {
	ast = &asset{
		path: path,
		name: filepath.ToSlash(name),
		fi:   fi,
	}

	if len(realPath) == 0 {
		ast.funcName = cfg.AssetPrefix + normalize(name)
	} else {
		ast.funcName = cfg.AssetPrefix + normalize(realPath)
	}
	return ast
}
