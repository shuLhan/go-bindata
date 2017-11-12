// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) (err error) {
	c.cwd, err = os.Getwd()
	if err != nil {
		return ErrCWD
	}

	// Ensure our configuration has sane values.
	err = c.validate()
	if err != nil {
		return
	}

	var toc []Asset
	var knownFuncs = make(map[string]int)
	var visitedPaths = make(map[string]bool)

	// Locate all the assets.
	for _, input := range c.Input {
		err = findFiles(c, input.Path, input.Recursive, &toc,
			knownFuncs, visitedPaths)
		if err != nil {
			return
		}
	}

	if c.Split {
		return translateToDir(c, toc)
	}

	return translateToFile(c, toc)
}

// ByName implement sort.Interface for []os.FileInfo based on Name()
type ByName []os.FileInfo

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

func getListFileInfo(
	dirpath string, visitedPaths map[string]bool,
) (
	list []os.FileInfo, newdirpath string, err error,
) {
	in, err := os.Stat(dirpath)
	if err != nil {
		return
	}

	if !in.IsDir() {
		newdirpath = filepath.Dir(dirpath)
		list = make([]os.FileInfo, 1)
		list[0] = in
		return
	}

	newdirpath = dirpath
	visitedPaths[dirpath] = true

	fd, err := os.Open(dirpath)
	if err != nil {
		_ = fd.Close()
		return
	}

	list, err = fd.Readdir(0)
	if err != nil {
		_ = fd.Close()
		return
	}

	err = fd.Close()
	if err != nil {
		return
	}

	// Sort to make output stable between invocations
	sort.Sort(ByName(list))

	return
}

func isIgnored(c *Config, path string) bool {
	for _, re := range c.Ignore {
		if re.MatchString(path) {
			return true
		}
	}

	for _, re := range c.Include {
		if re.MatchString(path) {
			return false
		}
	}

	return len(c.Include) > 0
}

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(
	c *Config,
	dir string,
	recursive bool,
	toc *[]Asset,
	knownFuncs map[string]int,
	visitedPaths map[string]bool,
) (err error) {
	var dirpath string

	if c.Prefix != nil {
		dirpath, err = filepath.Abs(dir)
		if err != nil {
			return
		}
	} else {
		dirpath = dir
	}

	list, dirpath, err := getListFileInfo(dirpath, visitedPaths)
	if err != nil {
		return
	}

	for _, file := range list {
		var asset Asset
		asset.Path = filepath.Join(dirpath, file.Name())
		asset.Name = filepath.ToSlash(asset.Path)

		if isIgnored(c, asset.Path) {
			continue
		}

		if file.IsDir() {
			if recursive {
				recursivePath := filepath.Join(dir, file.Name())
				visitedPaths[asset.Path] = true
				err = findFiles(c, recursivePath, recursive,
					toc, knownFuncs, visitedPaths)
				if err != nil {
					return
				}
			}
			continue
		}

		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			var linkPath string
			if linkPath, err = os.Readlink(asset.Path); err != nil {
				return
			}

			if !filepath.IsAbs(linkPath) {
				if linkPath, err = filepath.Abs(dirpath + "/" + linkPath); err != nil {
					return
				}
			}

			if _, ok := visitedPaths[linkPath]; !ok {
				visitedPaths[linkPath] = true
				err = findFiles(c, asset.Path, recursive, toc,
					knownFuncs, visitedPaths)
				if err != nil {
					return
				}
			}
			continue
		}

		if c.Prefix != nil && c.Prefix.MatchString(asset.Name) {
			asset.Name = c.Prefix.ReplaceAllString(asset.Name, "")
		} else if strings.HasSuffix(dir, file.Name()) {
			// Issue 110: dir is a full path, including
			// the file name (minus the basedir), so this
			// is what we have to use.
			asset.Name = dir
		} else {
			// Issue 110: dir is just that, a plain
			// directory, so we have to add the file's
			// name to it to form the full asset path.
			asset.Name = filepath.Join(dir, file.Name())
		}

		// If we have a leading slash, get rid of it.
		if len(asset.Name) > 0 && asset.Name[0] == '/' {
			asset.Name = asset.Name[1:]
		}

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		asset.Name = filepath.ToSlash(asset.Name)

		asset.Func = safeFunctionName(asset.Name, knownFuncs)
		asset.Path, _ = filepath.Abs(asset.Path)
		*toc = append(*toc, asset)
	}

	return
}

var regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// safeFunctionName converts the given name into a name
// which qualifies as a valid function identifier. It
// also compares against a known list of functions to
// prevent conflict based on name translation.
func safeFunctionName(name string, knownFuncs map[string]int) string {
	var inBytes, outBytes []byte
	var toUpper bool

	name = strings.ToLower(name)
	inBytes = []byte(name)

	for i := 0; i < len(inBytes); i++ {
		if regFuncName.Match([]byte{inBytes[i]}) {
			toUpper = true
		} else if toUpper {
			outBytes = append(outBytes, []byte(strings.ToUpper(string(inBytes[i])))...)
			toUpper = false
		} else {
			outBytes = append(outBytes, inBytes[i])
		}
	}

	name = string(outBytes)

	// Identifier can't start with a digit.
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}

	if num, ok := knownFuncs[name]; ok {
		knownFuncs[name] = num + 1
		name = fmt.Sprintf("%s%d", name, num)
	} else {
		knownFuncs[name] = 2
	}

	return name
}

func writeHeader(bfd io.Writer, c *Config, toc []Asset) (err error) {
	// Write the header. This makes e.g. Github ignore diffs in generated files.
	_, err = fmt.Fprint(bfd, headerGeneratedBy)
	if err != nil {
		return
	}

	if c.Split {
		_, err = fmt.Fprint(bfd, "// -- Common file --\n")
		if err != nil {
			return
		}
	} else {
		_, err = fmt.Fprint(bfd, "// sources:\n")
		if err != nil {
			return
		}

		for _, asset := range toc {
			relative, _ := filepath.Rel(c.cwd, asset.Path)

			_, err = fmt.Fprintf(bfd, "// %s\n", filepath.ToSlash(relative))
			if err != nil {
				return
			}
		}
	}

	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		if _, err = fmt.Fprintf(bfd, "// +build %s\n\n", c.Tags); err != nil {
			return
		}
	}

	return
}

//
// flushAndClose will flush the buffered writer `bfd` and close the file `fd`.
//
func flushAndClose(fd io.Closer, bfd *bufio.Writer, errParam error) (err error) {
	err = errParam

	if err == nil {
		err = bfd.Flush()
	}

	errClose := fd.Close()
	if errClose != nil {
		if err == nil {
			err = errClose
		}
	}

	return

}
