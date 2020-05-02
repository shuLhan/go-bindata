// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		cwd: cwd,
	}

	scanner := NewFSScanner(cfg)

	cases := []struct {
		desc      string
		inputs    []*InputConfig
		expError  error
		expAssets map[string]Asset
	}{{
		desc: "With single file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc/file1",
			Recursive: false,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
		},
	}, {
		desc: "With single directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With directory and a file",
		inputs: []*InputConfig{{
			Path: "./testdata/in/a",
		}, {
			Path: "./testdata/in/a/test.asset",
		}},
		expAssets: map[string]Asset{
			"testdata/in/a/test.asset": {
				Path: "testdata/in/a/test.asset",
				Name: "testdata/in/a/test.asset",
				Func: "bindataTestdataInATestasset",
			},
		},
	}, {
		desc: "With symlink to file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkFile/file1": {
				Path: "testdata/symlinkFile/file1",
				Name: "testdata/symlinkFile/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
		},
	}, {
		desc: "With symlink to file and duplicate",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkFile/file1": {
				Path: "testdata/symlinkFile/file1",
				Name: "testdata/symlinkFile/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
		},
	}, {
		desc: "With symlink to file and duplicate (reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkFile/file1": {
				Path: "testdata/symlinkFile/file1",
				Name: "testdata/symlinkFile/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With symlink to parent directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkParent",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkParent/symlinkTarget/file1": {
				Path: "testdata/symlinkParent/symlinkTarget/file1",
				Name: "testdata/symlinkParent/symlinkTarget/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkParent/symlinkTarget/file2": {
				Path: "testdata/symlinkParent/symlinkTarget/file2",
				Name: "testdata/symlinkParent/symlinkTarget/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkParent/symlinkTarget/file3": {
				Path: "testdata/symlinkParent/symlinkTarget/file3",
				Name: "testdata/symlinkParent/symlinkTarget/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkParent/symlinkTarget/file4": {
				Path: "testdata/symlinkParent/symlinkTarget/file4",
				Name: "testdata/symlinkParent/symlinkTarget/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With symlink to parent directory (in reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkParent",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkParent/symlinkTarget/file1": {
				Path: "testdata/symlinkParent/symlinkTarget/file1",
				Name: "testdata/symlinkParent/symlinkTarget/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkParent/symlinkTarget/file2": {
				Path: "testdata/symlinkParent/symlinkTarget/file2",
				Name: "testdata/symlinkParent/symlinkTarget/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkParent/symlinkTarget/file3": {
				Path: "testdata/symlinkParent/symlinkTarget/file3",
				Name: "testdata/symlinkParent/symlinkTarget/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkParent/symlinkTarget/file4": {
				Path: "testdata/symlinkParent/symlinkTarget/file4",
				Name: "testdata/symlinkParent/symlinkTarget/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With recursive symlink to directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkRecursiveParent",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkRecursiveParent/file1": {
				Path: "testdata/symlinkRecursiveParent/file1",
				Name: "testdata/symlinkRecursiveParent/file1",
				Func: "bindataTestdataSymlinkRecursiveParentFile1",
			},
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With recursive symlink to directory (in reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkRecursiveParent",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/symlinkSrc/file1": {
				Path: "testdata/symlinkSrc/file1",
				Name: "testdata/symlinkSrc/file1",
				Func: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				Path: "testdata/symlinkSrc/file2",
				Name: "testdata/symlinkSrc/file2",
				Func: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				Path: "testdata/symlinkSrc/file3",
				Name: "testdata/symlinkSrc/file3",
				Func: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				Path: "testdata/symlinkSrc/file4",
				Name: "testdata/symlinkSrc/file4",
				Func: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkRecursiveParent/file1": {
				Path: "testdata/symlinkRecursiveParent/file1",
				Name: "testdata/symlinkRecursiveParent/file1",
				Func: "bindataTestdataSymlinkRecursiveParentFile1",
			},
		},
	}, {
		desc: "With false duplicate function name",
		inputs: []*InputConfig{{
			Path:      "./testdata/dupname",
			Recursive: true,
		}},
		expAssets: map[string]Asset{
			"testdata/dupname/foo/bar": {
				Path: "testdata/dupname/foo/bar",
				Name: "testdata/dupname/foo/bar",
				Func: "bindataTestdataDupnameFooBar",
			},
			"testdata/dupname/foo_bar": {
				Path: "testdata/dupname/foo_bar",
				Name: "testdata/dupname/foo_bar",
				Func: "bindataTestdataDupnameFoobar",
			}},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		scanner.Reset()

		assets := make(map[string]Asset, len(c.expAssets))

		for _, in := range c.inputs {
			err = scanner.Scan(in.Path, "", in.Recursive)
			if err != nil {
				assert(t, c.expError, err, true)
				break
			}

			for k, asset := range scanner.assets {
				_, ok := assets[k]
				if !ok {
					assets[k] = asset
				}
			}

			scanner.Reset()
		}

		assert(t, len(c.expAssets), len(assets), true)

		for name, gotAsset := range assets {
			gotAsset.fi = nil
			assert(t, c.expAssets[name], gotAsset, true)
		}
	}
}

func TestScanAbsoluteSymlink(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		cwd: cwd,
	}

	tmpDir, err := ioutil.TempDir("", "go-bindata-test")
	if err != nil {
		t.Fatal(err)
	}
	//defer os.Remove(tmpDir)

	link := filepath.Join(tmpDir, "file1")
	target, err := filepath.Abs(filepath.Join("testdata", "symlinkSrc", "file1"))
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink(target, link)
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewFSScanner(cfg)
	err = scanner.Scan(tmpDir, "", true)
	if err != nil {
		t.Fatal(err)
	}

	if len(scanner.assets) != 1 {
		t.Fatalf("Expected exactly 1 asset but got %d", len(scanner.assets))
	}

	var actual Asset
	for _, asset := range scanner.assets {
		actual = asset
	}

	assert(t, link, actual.Path, true)
	assert(t, link, actual.Name, true)
}
