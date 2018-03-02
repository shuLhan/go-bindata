// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain
// Dedication license.  Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

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
		expAssets []*Asset
	}{{
		desc: "With single file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc/file1",
			Recursive: false,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}},
	}, {
		desc: "With single directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}},
	}, {
		desc: "With symlink to file",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkFile/file1",
			Name: "testdata/symlinkFile/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}},
	}, {
		desc: "With symlink to file and duplicate",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}, {
			Path: "testdata/symlinkFile/file1",
			Name: "testdata/symlinkFile/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}},
	}, {
		desc: "With symlink to file and duplicate (reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkFile",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkFile/file1",
			Name: "testdata/symlinkFile/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}},
	}, {
		desc: "With symlink to parent directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkParent",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkParent/symlinkTarget/file1",
			Name: "testdata/symlinkParent/symlinkTarget/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file2",
			Name: "testdata/symlinkParent/symlinkTarget/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file3",
			Name: "testdata/symlinkParent/symlinkTarget/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file4",
			Name: "testdata/symlinkParent/symlinkTarget/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}, {
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}},
	}, {
		desc: "With symlink to parent directory (in reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkParent",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file1",
			Name: "testdata/symlinkParent/symlinkTarget/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file2",
			Name: "testdata/symlinkParent/symlinkTarget/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file3",
			Name: "testdata/symlinkParent/symlinkTarget/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkParent/symlinkTarget/file4",
			Name: "testdata/symlinkParent/symlinkTarget/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}},
	}, {
		desc: "With recursive symlink to directory",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkRecursiveParent",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkRecursiveParent/file1",
			Name: "testdata/symlinkRecursiveParent/file1",
			Func: "bindataTestdataSymlinkRecursiveParentFile1",
		}, {
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}},
	}, {
		desc: "With recursive symlink to directory (in reverse order)",
		inputs: []*InputConfig{{
			Path:      "./testdata/symlinkSrc",
			Recursive: true,
		}, {
			Path:      "./testdata/symlinkRecursiveParent",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/symlinkSrc/file1",
			Name: "testdata/symlinkSrc/file1",
			Func: "bindataTestdataSymlinkSrcFile1",
		}, {
			Path: "testdata/symlinkSrc/file2",
			Name: "testdata/symlinkSrc/file2",
			Func: "bindataTestdataSymlinkSrcFile2",
		}, {
			Path: "testdata/symlinkSrc/file3",
			Name: "testdata/symlinkSrc/file3",
			Func: "bindataTestdataSymlinkSrcFile3",
		}, {
			Path: "testdata/symlinkSrc/file4",
			Name: "testdata/symlinkSrc/file4",
			Func: "bindataTestdataSymlinkSrcFile4",
		}, {
			Path: "testdata/symlinkRecursiveParent/file1",
			Name: "testdata/symlinkRecursiveParent/file1",
			Func: "bindataTestdataSymlinkRecursiveParentFile1",
		}},
	}, {
		desc: "With false duplicate function name",
		inputs: []*InputConfig{{
			Path:      "./testdata/dupname",
			Recursive: true,
		}},
		expAssets: []*Asset{{
			Path: "testdata/dupname/foo/bar",
			Name: "testdata/dupname/foo/bar",
			Func: "bindataTestdataDupnameFooBar",
		}, {
			Path: "testdata/dupname/foo_bar",
			Name: "testdata/dupname/foo_bar",
			Func: "bindataTestdataDupnameFoobar",
		}},
	}}

	for _, c := range cases {
		t.Log(c.desc)

		scanner.Reset()

		assets := make([]Asset, 0)

		for _, in := range c.inputs {
			err = scanner.Scan(in.Path, "", in.Recursive)
			if err != nil {
				assert(t, c.expError, err, true)
			}

			assets = append(assets, scanner.assets...)

			scanner.Reset()
		}

		assert(t, len(c.expAssets), len(assets), true)

		for x, gotAsset := range assets {
			assert(t, c.expAssets[x].Path, gotAsset.Path, true)
			assert(t, c.expAssets[x].Name, gotAsset.Name, true)
			assert(t, c.expAssets[x].Func, gotAsset.Func, true)
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

	actual := scanner.assets[0]

	assert(t, link, actual.Path, true)
	assert(t, link, actual.Name, true)
}
