package bindata

import (
	"testing"
)

func TestAsset(t *testing.T) {
	tests := []struct {
		desc       string
		assetName  string
		expContent string
		expErr     string
	}{{
		desc:       `With asset "in/split/test.1"`,
		assetName:  `in/split/test.1`,
		expContent: "// sample file 1\n",
	}, {
		desc:       `With asset "in/split/test.2"`,
		assetName:  `in/split/test.2`,
		expContent: "// sample file 2\n",
	}, {
		desc:       `With asset "in/split/test.3"`,
		assetName:  `in/split/test.3`,
		expContent: "// sample file 3\n",
	}, {
		desc:       `With asset "in/split/test.4"`,
		assetName:  `in/split/test.4`,
		expContent: "// sample file 4\n",
	}, {
		desc:      `With non existing asset "in/split/test.5"`,
		assetName: `in/split/test.5`,
		expErr:    "open in/split/test.5: file does not exist",
	}}

	for _, test := range tests {
		t.Log(test.desc)

		got, err := Asset(test.assetName)
		if err != nil {
			assert(t, test.expErr, err.Error(), true)
			continue
		}

		assert(t, test.expContent, string(got), true)
	}
}
