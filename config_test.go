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
		expErr: `Failed to stat input path '': lstat : no such file or directory`,
	}, {
		desc: `With directory not exist`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "./notexist",
			}},
		},
		expErr: `Failed to stat input path './notexist': lstat ./notexist: no such file or directory`,
	}, {
		desc: `With file as input`,
		cfg: &Config{
			Input: []InputConfig{{
				Path: "./README.md",
			}},
		},
	}}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.cfg.validateInput()
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}
	}
}

func TestValidateOutput(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		desc      string
		cfg       *Config
		expErr    string
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
		expErr: `Create output directory: mkdir /root/.ssh/: permission denied`,
	}, {
		desc: `With unwriteable file`,
		cfg: &Config{
			Output: "/template.go",
		},
		expErr: `open /template.go: permission denied`,
	}, {
		desc: `With output as directory`,
		cfg: &Config{
			Output: "/tmp/",
		},
		expOutput: filepath.Join("/tmp", DefOutputName),
	}}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.cfg.validateOutput()
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.expOutput, test.cfg.Output, true)
	}
}
