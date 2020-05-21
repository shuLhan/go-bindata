// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// byName implement sort.Interface for []os.FileInfo based on Name()
type byName []os.FileInfo

func (v byName) Len() int           { return len(v) }
func (v byName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

//
// fsScanner implement the file system scanner.
//
type fsScanner struct {
	cfg         *Config
	knownFuncs  map[string]int
	visitedDirs map[string]bool
	assets      map[string]*asset
	depth       int
}

//
// newFSScanner will create and initialize new file system scanner.
//
func newFSScanner(cfg *Config) (fss *fsScanner) {
	fss = &fsScanner{
		cfg: cfg,
	}

	fss.Reset()

	return
}

//
// Reset will clear all previous mapping and assets.
//
func (fss *fsScanner) Reset() {
	fss.knownFuncs = make(map[string]int)
	fss.visitedDirs = make(map[string]bool)
	fss.assets = make(map[string]*asset, 0)
	fss.depth = 0
}

//
// isIgnored will return,
// (1) true, if `path` is matched with one of ignore-pattern,
// (2) false, if `path` is matched with one of include-pattern,
// (3) true, if include-pattern is defined but no matched found.
//
func (fss *fsScanner) isIgnored(path string) bool {
	// (1)
	for _, re := range fss.cfg.Ignore {
		if re.MatchString(path) {
			return true
		}
	}

	// (2)
	for _, re := range fss.cfg.Include {
		if re.MatchString(path) {
			return false
		}
	}

	// (3)
	return len(fss.cfg.Include) > 0
}

func (fss *fsScanner) cleanPrefix(path string) string {
	if fss.cfg.Prefix == nil {
		return path
	}

	return fss.cfg.Prefix.ReplaceAllString(path, "")
}

//
// addAsset will add new asset based on path, realPath, and file info.  The
// path can be a directory or file. Realpath reference to the original path if
// path is symlink, if path is not symlink then path and realPath will be equal.
//
func (fss *fsScanner) addAsset(path, realPath string, fi os.FileInfo) {
	name := fss.cleanPrefix(path)

	asset := newAsset(fss.cfg, path, name, realPath, fi)

	// Check if the asset's name is already exist.
	_, ok := fss.assets[name]
	if ok {
		if fss.cfg.Verbose {
			fmt.Printf("= %+v\n", path)
		}
		return
	}

	num, ok := fss.knownFuncs[asset.funcName]
	if ok {
		fss.knownFuncs[asset.funcName] = num + 1
		asset.funcName = fmt.Sprintf("%s_%d", asset.funcName, num)
	} else {
		fss.knownFuncs[asset.funcName] = 2
	}

	if fss.cfg.Verbose {
		fmt.Printf("+ %+v\n", path)
	}

	fss.assets[name] = asset
}

//
// getListFileInfo will return list of files in `path`.
//
// (1) set the path visited status to true,
// (2) read all files inside directory into list, and
// (3) sort the list to make output stable between invocations.
//
func (fss *fsScanner) getListFileInfo(path string) (
	list []os.FileInfo, err error,
) {
	// (1)
	fss.visitedDirs[path] = true

	// (2)
	fd, err := os.Open(path)
	if err != nil {
		_ = fd.Close()
		return
	}

	list, err = fd.Readdir(0)
	if err != nil {
		_ = fd.Close()
		return
	}

	err = fd.Close()
	if err != nil {
		return
	}

	// (3)
	sort.Sort(byName(list))

	return
}

//
// scanSymlink reads the real-path from symbolic link file and converts the path
// and real-path into a relative path.
//
func (fss *fsScanner) scanSymlink(path string, recursive bool) (
	err error,
) {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	fi, err := os.Lstat(realPath)
	if err != nil {
		return err
	}

	if fi.Mode().IsRegular() {
		fss.addAsset(path, realPath, fi)
		return nil
	}

	if !recursive {
		if fss.depth > 0 {
			return nil
		}
		fss.depth++
	}

	_, ok := fss.visitedDirs[realPath]
	if ok {
		return nil
	}

	list, err := fss.getListFileInfo(path)
	if err != nil {
		return err
	}

	for _, fi = range list {
		filePath := filepath.Join(path, fi.Name())
		fileRealPath := filepath.Join(realPath, fi.Name())

		err = fss.Scan(filePath, fileRealPath, recursive)
		if err != nil {
			return err
		}
	}

	return nil
}

//
// Scan will scan the file or content of directory in `path`.
//
func (fss *fsScanner) Scan(path, realPath string, recursive bool) (err error) {
	path = filepath.Clean(path)

	if fss.isIgnored(path) {
		if fss.cfg.Verbose {
			fmt.Printf("- %s\n", path)
		}
		return nil
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return fss.scanSymlink(path, recursive)
	}

	if fi.Mode().IsRegular() {
		fss.addAsset(path, realPath, fi)
		return nil
	}

	if !recursive {
		if fss.depth > 0 {
			return nil
		}
		fss.depth++
	}

	_, ok := fss.visitedDirs[path]
	if ok {
		return nil
	}

	list, err := fss.getListFileInfo(path)
	if err != nil {
		return err
	}

	for _, fi = range list {
		filePath := filepath.Join(path, fi.Name())

		err = fss.Scan(filePath, "", recursive)
		if err != nil {
			return err
		}
	}

	return nil
}
