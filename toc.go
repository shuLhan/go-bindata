// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type assetTree struct {
	Asset    Asset
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

func (root *assetTree) Add(route []string, asset Asset) {
	for _, name := range route {
		root = root.child(name)
	}
	root.Asset = asset
}

func ident(w io.Writer, n int) (err error) {
	for i := 0; i < n; i++ {
		_, err = w.Write([]byte{'\t'})
		if err != nil {
			return
		}
	}
	return
}

func (root *assetTree) funcOrNil() string {
	if root.Asset.Func == "" {
		return "nil"
	}

	return root.Asset.Func
}

func (root *assetTree) writeGoMap(w io.Writer, nident int) (err error) {
	if nident == 0 {
		_, err = io.WriteString(w, "&bintree")
		if err != nil {
			return
		}
	}

	fmt.Fprintf(w, "{%s, map[string]*bintree{", root.funcOrNil())

	if len(root.Children) > 0 {
		_, err = io.WriteString(w, "\n")
		if err != nil {
			return
		}

		// Sort to make output stable between invocations
		filenames := make([]string, len(root.Children))
		i := 0
		for filename := range root.Children {
			filenames[i] = filename
			i++
		}
		sort.Strings(filenames)

		for _, p := range filenames {
			err = ident(w, nident+1)
			if err != nil {
				return
			}
			fmt.Fprintf(w, `"%s": `, p)

			err = root.Children[p].writeGoMap(w, nident+1)
			if err != nil {
				return
			}
		}

		err = ident(w, nident)
		if err != nil {
			return
		}
	}

	_, err = io.WriteString(w, "}}")
	if err != nil {
		return
	}

	if nident > 0 {
		_, err = io.WriteString(w, ",")
		if err != nil {
			return
		}
	}

	_, err = io.WriteString(w, "\n")

	return
}

func (root *assetTree) WriteAsGoMap(w io.Writer) (err error) {
	_, err = fmt.Fprint(w, `type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = `)
	if err != nil {
		return
	}

	return root.writeGoMap(w, 0)
}

func writeTOCTree(w io.Writer, toc []Asset) error {
	_, err := fmt.Fprintf(w, `// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

`)
	if err != nil {
		return err
	}
	tree := newAssetTree()
	for i := range toc {
		pathList := strings.Split(toc[i].Name, "/")
		tree.Add(pathList, toc[i])
	}

	return tree.WriteAsGoMap(w)
}

//
// getLongestAssetNameLen will return length of the longest asset name in toc.
//
func getLongestAssetNameLen(toc []Asset) (longest int) {
	for _, asset := range toc {
		lenName := len(asset.Name)
		if lenName > longest {
			longest = lenName
		}
	}

	return
}

// writeTOC writes the table of contents file.
func writeTOC(w io.Writer, toc []Asset) error {
	err := writeTOCHeader(w)
	if err != nil {
		return err
	}

	longestNameLen := getLongestAssetNameLen(toc)

	for i := range toc {
		err = writeTOCAsset(w, &toc[i], longestNameLen)
		if err != nil {
			return err
		}
	}

	return writeTOCFooter(w)
}

// writeTOCHeader writes the table of contents file header.
func writeTOCHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %%s can't read by error: %%v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %%s can't read by error: %%v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// AssetNames returns the names of the assets.
// nolint: deadcode
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
`)
	return err
}

// writeTOCAsset write a TOC entry for the given asset.
func writeTOCAsset(w io.Writer, asset *Asset, longestNameLen int) (err error) {
	toWrite := " "

	for x := 0; x < longestNameLen-len(asset.Name); x++ {
		toWrite += " "
	}

	toWrite = "\t\"" + asset.Name + "\":" + toWrite + asset.Func + ",\n"

	_, err = io.WriteString(w, toWrite)

	return
}

// writeTOCFooter writes the table of contents file footer.
func writeTOCFooter(w io.Writer) error {
	_, err := fmt.Fprintf(w, `}

`)
	return err
}
