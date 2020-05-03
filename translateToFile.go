// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"bufio"
	"fmt"
	"os"
)

// translateToFile generates one single file
func translateToFile(c *Config, keys []string, toc map[string]*asset) (err error) {
	// Create output file.
	fd, err := os.Create(c.Output)
	if err != nil {
		return err
	}

	if c.Verbose {
		fmt.Printf("> %s\n", c.Output)
	}

	// Create a buffered writer for better performance.
	bfd := bufio.NewWriter(fd)

	err = writeHeader(bfd, c, keys, toc)
	if err != nil {
		goto out
	}

	// Write package declaration.
	_, err = fmt.Fprintf(bfd, "\npackage %s\n\n", c.Package)
	if err != nil {
		goto out
	}

	// Write assets.
	if c.Debug || c.Dev {
		err = writeDebug(bfd, c, keys, toc)
	} else {
		err = writeRelease(bfd, c, keys, toc)
	}

	if err != nil {
		goto out
	}

	// Write table of contents
	err = writeTOC(bfd, keys, toc)
	if err != nil {
		goto out
	}

	// Write hierarchical tree of assets
	err = writeTOCTree(bfd, keys, toc)
	if err != nil {
		return err
	}

	// Write restore procedure
	err = writeRestore(bfd)
out:
	return flushAndClose(fd, bfd, err)
}
