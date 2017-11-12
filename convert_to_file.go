package bindata

import (
	"bufio"
	"fmt"
	"os"
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

	err = writeHeader(bfd, c, toc, wd)
	if err != nil {
		return err
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
