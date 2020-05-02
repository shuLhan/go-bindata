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
		expAssets map[string]*asset
	}{{
		desc: "With single file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc/file1",
			Recursive: false,
		}},
		expAssets: map[string]*asset{
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
		},
	}, {
		desc: "With single directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: map[string]*asset{
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
			},
		},
	}, {
		desc: "With directory and a file",
		inputs: []*InputConfig{{
			Path: "./testdata/in/a",
		}, {
			Path: "./testdata/in/a/test.asset",
		}},
		expAssets: map[string]*asset{
			"testdata/in/a/test.asset": {
				path:     "testdata/in/a/test.asset",
				name:     "testdata/in/a/test.asset",
				funcName: "bindataTestdataInATestasset",
			},
		},
	}, {
		desc: "With symlink to file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}},
		expAssets: map[string]*asset{
			"testdata/symlinkFile/file1": {
				path:     "testdata/symlinkFile/file1",
				name:     "testdata/symlinkFile/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
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
		expAssets: map[string]*asset{
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkFile/file1": {
				path:     "testdata/symlinkFile/file1",
				name:     "testdata/symlinkFile/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
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
		expAssets: map[string]*asset{
			"testdata/symlinkFile/file1": {
				path:     "testdata/symlinkFile/file1",
				name:     "testdata/symlinkFile/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
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
		expAssets: map[string]*asset{
			"testdata/symlinkParent/symlinkTarget/file1": {
				path:     "testdata/symlinkParent/symlinkTarget/file1",
				name:     "testdata/symlinkParent/symlinkTarget/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkParent/symlinkTarget/file2": {
				path:     "testdata/symlinkParent/symlinkTarget/file2",
				name:     "testdata/symlinkParent/symlinkTarget/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkParent/symlinkTarget/file3": {
				path:     "testdata/symlinkParent/symlinkTarget/file3",
				name:     "testdata/symlinkParent/symlinkTarget/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkParent/symlinkTarget/file4": {
				path:     "testdata/symlinkParent/symlinkTarget/file4",
				name:     "testdata/symlinkParent/symlinkTarget/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
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
		expAssets: map[string]*asset{
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkParent/symlinkTarget/file1": {
				path:     "testdata/symlinkParent/symlinkTarget/file1",
				name:     "testdata/symlinkParent/symlinkTarget/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkParent/symlinkTarget/file2": {
				path:     "testdata/symlinkParent/symlinkTarget/file2",
				name:     "testdata/symlinkParent/symlinkTarget/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkParent/symlinkTarget/file3": {
				path:     "testdata/symlinkParent/symlinkTarget/file3",
				name:     "testdata/symlinkParent/symlinkTarget/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkParent/symlinkTarget/file4": {
				path:     "testdata/symlinkParent/symlinkTarget/file4",
				name:     "testdata/symlinkParent/symlinkTarget/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
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
		expAssets: map[string]*asset{
			"testdata/symlinkRecursiveParent/file1": {
				path:     "testdata/symlinkRecursiveParent/file1",
				name:     "testdata/symlinkRecursiveParent/file1",
				funcName: "bindataTestdataSymlinkRecursiveParentFile1",
			},
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
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
		expAssets: map[string]*asset{
			"testdata/symlinkSrc/file1": {
				path:     "testdata/symlinkSrc/file1",
				name:     "testdata/symlinkSrc/file1",
				funcName: "bindataTestdataSymlinkSrcFile1",
			},
			"testdata/symlinkSrc/file2": {
				path:     "testdata/symlinkSrc/file2",
				name:     "testdata/symlinkSrc/file2",
				funcName: "bindataTestdataSymlinkSrcFile2",
			},
			"testdata/symlinkSrc/file3": {
				path:     "testdata/symlinkSrc/file3",
				name:     "testdata/symlinkSrc/file3",
				funcName: "bindataTestdataSymlinkSrcFile3",
			},
			"testdata/symlinkSrc/file4": {
				path:     "testdata/symlinkSrc/file4",
				name:     "testdata/symlinkSrc/file4",
				funcName: "bindataTestdataSymlinkSrcFile4",
			},
			"testdata/symlinkRecursiveParent/file1": {
				path:     "testdata/symlinkRecursiveParent/file1",
				name:     "testdata/symlinkRecursiveParent/file1",
				funcName: "bindataTestdataSymlinkRecursiveParentFile1",
			},
		},
	}, {
		desc: "With false duplicate function name",
		inputs: []*InputConfig{{
			Path:      "./testdata/dupname",
			Recursive: true,
		}},
		expAssets: map[string]*asset{
			"testdata/dupname/foo/bar": {
				path:     "testdata/dupname/foo/bar",
				name:     "testdata/dupname/foo/bar",
				funcName: "bindataTestdataDupnameFooBar",
			},
			"testdata/dupname/foo_bar": {
				path:     "testdata/dupname/foo_bar",
				name:     "testdata/dupname/foo_bar",
				funcName: "bindataTestdataDupnameFoobar",
			}},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		scanner.Reset()

		assets := make(map[string]*asset, len(c.expAssets))

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

	var actual *asset
	for _, asset := range scanner.assets {
		actual = asset
	}

	assert(t, link, actual.path, true)
	assert(t, link, actual.name, true)
}
