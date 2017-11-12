package bindata

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// translateToDir generates splited file
func translateToDir(c *Config, toc []Asset, wd string) error {
	if err := generateCommonFile(c, toc, wd); err != nil {
		return err
	}

	for i := range toc {
		if err := generateOneAsset(c, &toc[i], wd); err != nil {
			return err
		}
	}

	return nil
}

func generateCommonFile(c *Config, toc []Asset, wd string) error {
	// Create output file.
	fd, err := os.Create(filepath.Join(c.Output, DefOutputName))
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
		err = writeDebugHeader(bfd)
	} else {
		err = writeReleaseHeader(bfd, c)
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

func generateOneAsset(c *Config, a *Asset, wd string) error {
	// Create output file.
	fd, err := os.Create(filepath.Join(c.Output, a.Func+".go"))
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

	if _, err = fmt.Fprint(bfd, "// source: "); err != nil {
		return err
	}

	relative, _ := filepath.Rel(wd, a.Path)
	if _, err = fmt.Fprintln(bfd, filepath.ToSlash(relative)); err != nil {
		return err
	}

	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		if _, err = fmt.Fprintf(bfd, "// +build %s\n\n", c.Tags); err != nil {
			return err
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(bfd, "package %s\n\n", c.Package)
	if err != nil {
		return err
	}

	// Write assets.
	if c.Debug || c.Dev {
		err = writeOneFileDebug(bfd, c, a)
	} else {
		err = writeOneFileRelease(bfd, c, a)
	}

	return err
}
