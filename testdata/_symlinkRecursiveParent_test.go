package bindata

import (
	"testing"
)

func TestAsset(t *testing.T) {
	tests := []struct {
		desc   string
		name   string
		exp    string
		expErr string
	}{{
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent",
		expErr: "open symlinkRecursiveParent: file does not exist",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent/symlinkTarget",
		expErr: "open symlinkRecursiveParent/symlinkTarget: file does not exist",
	}, {
		desc: "With valid asset",
		name: "symlinkRecursiveParent/file1",
		exp:  "// symlinkRecursiveParent/file1\n",
	}, {
		desc:   "With invalid asset",
		name:   "symlinkRecursiveParent/symlinkTarget/file1",
		expErr: "open symlinkRecursiveParent/symlinkTarget/file1: file does not exist",
	}}

	for _, test := range tests {
		t.Log(test.desc, ":", test.name)

		got, err := Asset(test.name)
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.exp, string(got), true)
	}
}
