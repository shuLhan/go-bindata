// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type assetTree struct {
	asset    *asset
	Children map[string]*assetTree
}

func newAssetTree() *assetTree {
	tree := &assetTree{}
	tree.Children = make(map[string]*assetTree)
	return tree
}

func (root *assetTree) child(name string) *assetTree {
	rv, ok := root.Children[name]
	if !ok {
		rv = newAssetTree()
		root.Children[name] = rv
	}
	return rv
}

func (root *assetTree) Add(route []string, ast *asset) {
	for _, name := range route {
		root = root.child(name)
	}
	root.asset = ast
}

func ident(w io.Writer, n int) (err error) {
	for i := 0; i < n; i++ {
		_, err = w.Write([]byte{'\t'})
		if err != nil {
			return err
		}
	}
	return nil
}

func (root *assetTree) funcOrNil() string {
	if root.asset == nil || len(root.asset.funcName) == 0 {
		return "nil"
	}
	return root.asset.funcName
}

//
// getFilenames will return all files sorted, to make output stable between
// invocations.
//
func (root *assetTree) getFilenames() (filenames []string) {
	filenames = make([]string, len(root.Children))
	x := 0
	for filename := range root.Children {
		filenames[x] = filename
		x++
	}
	sort.Strings(filenames)

	return
}

func (root *assetTree) writeGoMap(w io.Writer, nident int) (err error) {
	fmt.Fprintf(w, tmplBinTreeValues, root.funcOrNil())

	if len(root.Children) > 0 {
		_, err = io.WriteString(w, "\n")
		if err != nil {
			return err
		}

		filenames := root.getFilenames()

		for _, p := range filenames {
			err = ident(w, nident+1)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, `"%s": `, p)

			err = root.Children[p].writeGoMap(w, nident+1)
			if err != nil {
				return err
			}
		}

		err = ident(w, nident)
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, "}}")
	if err != nil {
		return err
	}

	if nident > 0 {
		_, err = io.WriteString(w, ",")
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, "\n")

	return err
}

func (root *assetTree) WriteAsGoMap(w io.Writer) (err error) {
	_, err = fmt.Fprint(w, tmplTypeBintree)
	if err != nil {
		return
	}

	return root.writeGoMap(w, 0)
}

func writeTOCTree(w io.Writer, keys []string, toc map[string]*asset) error {
	_, err := fmt.Fprint(w, tmplFuncAssetDir)
	if err != nil {
		return err
	}
	tree := newAssetTree()
	for _, key := range keys {
		ast := toc[key]
		pathList := strings.Split(ast.name, "/")
		tree.Add(pathList, ast)
	}

	return tree.WriteAsGoMap(w)
}

//
// getLongestAssetNameLen will return length of the longest asset name in toc.
//
func getLongestAssetNameLen(keys []string) (longest int) {
	for _, key := range keys {
		lenName := len(key)
		if lenName > longest {
			longest = lenName
		}
	}
	return longest
}

// writeTOC writes the table of contents file.
func writeTOC(w io.Writer, keys []string, toc map[string]*asset) (err error) {
	_, err = fmt.Fprint(w, tmplFuncAsset)
	if err != nil {
		return err
	}

	longestNameLen := getLongestAssetNameLen(keys)

	for _, key := range keys {
		err = writeTOCAsset(w, toc[key], longestNameLen)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "}\n")

	return
}

// writeTOCAsset write a TOC entry for the given asset.
func writeTOCAsset(w io.Writer, ast *asset, longestNameLen int) (err error) {
	toWrite := " "

	for x := 0; x < longestNameLen-len(ast.name); x++ {
		toWrite += " "
	}

	toWrite = "\t\"" + ast.name + "\":" + toWrite + ast.funcName + ",\n"

	_, err = io.WriteString(w, toWrite)

	return
}
