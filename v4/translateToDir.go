// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// translateToDir generates splited file
func translateToDir(c *Config, keys []string, toc map[string]*asset) error {
	if err := generateCommonFile(c, keys, toc); err != nil {
		return err
	}

	for _, key := range keys {
		ast := toc[key]
		if err := generateOneAsset(c, ast); err != nil {
			return err
		}
	}

	return nil
}

func generateCommonFile(c *Config, keys []string, toc map[string]*asset) (err error) {
	// Create output file.
	out := filepath.Join(c.Output, DefOutputName)
	fd, err := os.Create(out)
	if err != nil {
		return err
	}

	if c.Verbose {
		fmt.Printf("> %s\n", out)
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
		err = writeDebugHeader(bfd)
	} else {
		err = writeReleaseHeader(bfd, c)
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
		goto out
	}

	// Write restore procedure
	err = writeRestore(bfd)

out:
	return flushAndClose(fd, bfd, err)
}

func generateOneAsset(c *Config, ast *asset) (err error) {
	// Create output file.
	out := filepath.Join(c.Output, ast.funcName+".go")
	fd, err := os.Create(out)
	if err != nil {
		return err
	}

	if c.Verbose {
		fmt.Printf("> %s\n", out)
	}

	// Create a buffered writer for better performance.
	bfd := bufio.NewWriter(fd)

	// Write the header. This makes e.g. Github ignore diffs in generated files.
	_, err = fmt.Fprint(bfd, headerGeneratedBy)
	if err != nil {
		goto out
	}

	if _, err = fmt.Fprint(bfd, "// source: "); err != nil {
		goto out
	}

	if _, err = fmt.Fprintln(bfd, ast.path); err != nil {
		goto out
	}

	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		if _, err = fmt.Fprintf(bfd, "// +build %s\n\n", c.Tags); err != nil {
			goto out
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(bfd, "package %s\n\n", c.Package)
	if err != nil {
		goto out
	}

	// Write assets.
	if c.Debug || c.Dev {
		err = writeOneFileDebug(bfd, c, ast)
	} else {
		err = writeOneFileRelease(bfd, c, ast)
	}
out:
	return flushAndClose(fd, bfd, err)
}
