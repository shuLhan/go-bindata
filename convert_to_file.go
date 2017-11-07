package bindata

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// translateToFile generates one single file
func translateToFile(c *Config, toc []Asset, wd string) error {
	// Create output file.
	fd, err := os.Create(c.Output)
	if err != nil {
		return err
	}

	defer fd.Close()

	// Create a buffered writer for better performance.
	bfd := bufio.NewWriter(fd)
	defer bfd.Flush()

	// Write the header. This makes e.g. Github ignore diffs in generated files.
	_, err = fmt.Fprint(bfd, headerGeneratedBy)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprint(bfd, "// sources:\n"); err != nil {
		return err
	}

	for _, asset := range toc {
		relative, _ := filepath.Rel(wd, asset.Path)
		if _, err = fmt.Fprintf(bfd, "// %s\n", filepath.ToSlash(relative)); err != nil {
			return err
		}
	}

	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		if _, err = fmt.Fprintf(bfd, "// +build %s\n\n", c.Tags); err != nil {
			return err
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(bfd, "\npackage %s\n\n", c.Package)
	if err != nil {
		return err
	}

	// Write assets.
	if c.Debug || c.Dev {
		err = writeDebug(bfd, c, toc)
	} else {
		err = writeRelease(bfd, c, toc)
	}

	if err != nil {
		return err
	}

	// Write table of contents
	if err := writeTOC(bfd, toc); err != nil {
		return err
	}
	// Write hierarchical tree of assets
	if err := writeTOCTree(bfd, toc); err != nil {
		return err
	}

	// Write restore procedure
	return writeRestore(bfd)
}
