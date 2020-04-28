// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		desc   string
		cfg    *Config
		exp    []InputConfig
		expErr string
	}{{
		desc:   `With empty list`,
		cfg:    &Config{},
		expErr: ErrNoInput.Error(),
	}, {
		desc: `With empty path`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "",
			}},
		},
		exp: []InputConfig{{
			Path: ".",
		}},
	}, {
		desc: `With directory not exist`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "./notexist",
			}},
		},
		expErr: `failed to stat input path 'notexist': lstat notexist: no such file or directory`,
	}, {
		desc: `With file as input`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "./README.md",
			}},
		},
		exp: []InputConfig{{
			Path: "README.md",
		}},
	}, {
		desc: `With duplicate inputs`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "./testdata/in/test.asset",
			}, {
				Path: "./testdata/in/test.asset",
			}},
		},
		exp: []InputConfig{{
			Path: "testdata/in/test.asset",
		}},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.cfg.validateInput()
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.exp, test.cfg.Input, true)
	}
}

func TestValidateOutput(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		desc      string
		cfg       *Config
		expErr    []string
		expOutput string
	}{{
		desc: `With empty`,
		cfg: &Config{
			cwd: cwd,
		},
		expOutput: filepath.Join(cwd, DefOutputName),
	}, {
		desc: `With unwriteable directory`,
		cfg: &Config{
			Output: "/root/.ssh/template.go",
		},
		expErr: []string{
			`create output directory: mkdir /root: read-only file system`,
			`create output directory: mkdir /root/.ssh/: permission denied`,
		},
	}, {
		desc: `With unwriteable file`,
		cfg: &Config{
			Output: "/template.go",
		},
		expErr: []string{
			`open /template.go: permission denied`,
			`open /template.go: read-only file system`,
		},
	}, {
		desc: `With output as directory`,
		cfg: &Config{
			Output: "/tmp/",
		},
		expOutput: filepath.Join("/tmp", DefOutputName),
	}}

test:
	for _, c := range cases {
		t.Log(c.desc)

		err := c.cfg.validateOutput()
		if err != nil {
			for _, expErr := range c.expErr {
				if expErr == err.Error() {
					continue test
				}
			}
			t.Fatalf("expecting one of error %q, got %q",
				c.expErr, err.Error())
		}

		assert(t, c.expOutput, c.cfg.Output, true)
	}
}
