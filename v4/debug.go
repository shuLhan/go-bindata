// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"fmt"
	"io"
	"path/filepath"
)

// writeOneFileDebug writes the debug code file for each file (when splited file).
func writeOneFileDebug(w io.Writer, c *Config, ast *asset) error {
	if err := writeDebugFileHeader(w, c.Dev); err != nil {
		return err
	}

	if err := writeDebugAsset(w, c, ast); err != nil {
		return err
	}

	return nil
}

// writeDebug writes the debug code file for single file.
func writeDebug(w io.Writer, c *Config, keys []string, toc map[string]*asset) error {
	err := writeDebugHeader(w)
	if err != nil {
		return err
	}

	for _, key := range keys {
		ast := toc[key]
		err = writeDebugAsset(w, c, ast)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeDebugHeader writes output file headers for each file.
// This targets debug builds.
func writeDebugFileHeader(w io.Writer, dev bool) error {
	add := ""
	if dev {
		add = `
	"path/filepath"`
	}

	_, err := fmt.Fprintf(w, `import (
	"fmt"
	"os"%s
)

`, add)

	return err
}

// writeDebugHeader writes output file headers for sigle file.
// This targets debug builds.
func writeDebugHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// bindataRead reads the given file from disk. It returns an error on failure.
func bindataRead(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %%s at %%s: %%v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

`)
	return err
}

// writeDebugAsset write a debug entry for the given asset.
// A debug entry is simply a function which reads the asset from
// the original file (e.g.: from disk).
func writeDebugAsset(w io.Writer, c *Config, ast *asset) error {
	pathExpr := fmt.Sprintf("%q", filepath.Join(c.cwd, ast.path))
	if c.Dev {
		pathExpr = fmt.Sprintf("filepath.Join(rootDir, %q)", ast.name)
	}

	_, err := fmt.Fprintf(w, `// %s reads file data from disk. It returns an error on failure.
func %sBytes() ([]byte, error) {
	asset, err := %s()
	if asset == nil {
		return nil, err
	}
	return asset.bytes, err
}

func %s() (*asset, error) {
	path := %s
	name := %q
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %%s at %%s: %%v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

`, ast.funcName, ast.funcName, ast.funcName, ast.funcName, pathExpr, ast.name)
	return err
}
