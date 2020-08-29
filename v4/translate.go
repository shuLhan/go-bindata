// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"sort"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) (err error) {
	// Ensure our configuration has sane values.
	err = c.validate()
	if err != nil {
		return
	}

	scanner := newFSScanner(c)

	assets := make(map[string]*asset, 0)

	// Locate all the assets.
	for _, input := range c.Input {
		err = scanner.Scan(input.Path, "", input.Recursive)
		if err != nil {
			return
		}

		for k, asset := range scanner.assets {
			_, ok := assets[k]
			if !ok {
				assets[k] = asset
			}
		}

		scanner.Reset()
	}

	keys := make([]string, 0, len(assets))
	for key := range assets {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if c.Split {
		return translateToDir(c, keys, assets)
	}

	return translateToFile(c, keys, assets)
}
