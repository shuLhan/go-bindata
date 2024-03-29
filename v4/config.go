// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	// DefPackageName define default package name.
	DefPackageName = "main"

	// DefOutputName define default generated file name.
	DefOutputName = "bindata.go"

	// Default prefix for asset functions
	DefAssetPrefixName = "bindata"
)

// List of errors.
var (
	ErrNoInput       = errors.New("no input")
	ErrNoPackageName = errors.New("missing package name")
	ErrCWD           = errors.New("unable to determine current working directory")
)

// Config defines a set of options for the asset conversion.
type Config struct {
	// cwd contains current working directory.
	cwd string

	// Name of the package to use. Defaults to 'main'.
	Package string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Input defines the directory path, containing all asset files as
	// well as whether to recursively process assets in any sub directories.
	Input []InputConfig

	// Output defines the output file for the generated code.
	// If left empty, this defaults to 'bindata.go' in the current
	// working directory and the current directory in case of having true
	// to `Split` config.
	Output string

	// This defines the string that is prepended to asset functions.
	// This can be used to export these functions directly.
	AssetPrefix string

	// Prefix defines a regular expression which should used to strip
	// substrings from all file names when generating the keys in the table of
	// contents.  For example, running without the `-prefix` flag, we get:
	//
	// 	$ go-bindata /path/to/templates
	// 	go_bindata["/path/to/templates/foo.html"] = _path_to_templates_foo_html
	//
	// Running with the `-prefix` flag, we get:
	//
	//	$ go-bindata -prefix "/.*/some/" /a/path/to/some/templates/
	//	_bindata["templates/foo.html"] = templates_foo_html
	Prefix *regexp.Regexp

	// Ignores any filenames matching the regex pattern specified, e.g.
	// path/to/file.ext will ignore only that file, or \\.gitignore
	// will match any .gitignore file.
	//
	// This parameter can be provided multiple times.
	Ignore []*regexp.Regexp

	// Include contains list of regex to filter input files.
	Include []*regexp.Regexp

	// When nonzero, use this as mode for all files.
	Mode uint

	// When nonzero, use this as unix timestamp for all files.
	ModTime int64

	// When true, size, mode and modtime are not preserved from files
	NoMetadata bool

	// NoMemCopy will alter the way the output file is generated.
	//
	// It will employ a hack that allows us to read the file data directly from
	// the compiled program's `.rodata` section. This ensures that when we call
	// call our generated function, we omit unnecessary mem copies.
	//
	// The downside of this, is that it requires dependencies on the `reflect` and
	// `unsafe` packages. These may be restricted on platforms like AppEngine and
	// thus prevent you from using this mode.
	//
	// Another disadvantage is that the byte slice we create, is strictly read-only.
	// For most use-cases this is not a problem, but if you ever try to alter the
	// returned byte slice, a runtime panic is thrown. Use this mode only on target
	// platforms where memory constraints are an issue.
	//
	// The default behaviour is to use the old code generation method. This
	// prevents the two previously mentioned issues, but will employ at least one
	// extra memcopy and thus increase memory requirements.
	//
	// For instance, consider the following two examples:
	//
	// This would be the default mode, using an extra memcopy but gives a safe
	// implementation without dependencies on `reflect` and `unsafe`:
	//
	// 	func myfile() []byte {
	// 		return []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a}
	// 	}
	//
	// Here is the same functionality, but uses the `.rodata` hack.
	// The byte slice returned from this example can not be written to without
	// generating a runtime error.
	//
	// 	var _myfile = "\x89\x50\x4e\x47\x0d\x0a\x1a"
	//
	// 	func myfile() []byte {
	// 		var empty [0]byte
	// 		sx := (*reflect.StringHeader)(unsafe.Pointer(&_myfile))
	// 		b := empty[:]
	// 		bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	// 		bx.Data = sx.Data
	// 		bx.Len = len(_myfile)
	// 		bx.Cap = bx.Len
	// 		return b
	// 	}
	NoMemCopy bool

	// NoCompress means the assets are /not/ GZIP compressed before being turned
	// into Go code. The generated function will automatically unzip
	// the file data when called. Defaults to false.
	NoCompress bool

	// Perform a debug build. This generates an asset file, which
	// loads the asset contents directly from disk at their original
	// location, instead of embedding the contents in the code.
	//
	// This is mostly useful if you anticipate that the assets are
	// going to change during your development cycle. You will always
	// want your code to access the latest version of the asset.
	// Only in release mode, will the assets actually be embedded
	// in the code. The default behaviour is Release mode.
	Debug bool

	// Perform a dev build, which is nearly identical to the debug option. The
	// only difference is that instead of absolute file paths in generated code,
	// it expects a variable, `rootDir`, to be set in the generated code's
	// package (the author needs to do this manually), which it then prepends to
	// an asset's name to construct the file path on disk.
	//
	// This is mainly so you can push the generated code file to a shared
	// repository.
	Dev bool

	// Split the output into several files. Every embedded file is bound into
	// a specific file, and a common file is also generated containing API and
	// other common parts.
	// If true, the output config is a directory and not a file.
	Split bool

	// MD5Checksum is a flag that, when set to true, indicates to calculate
	// MD5 checksums for files.
	MD5Checksum bool

	// Verbose flag to display verbose output.
	Verbose bool
}

// NewConfig returns a default configuration struct.
func NewConfig() *Config {
	c := new(Config)
	c.Package = DefPackageName
	c.Output = DefOutputName
	c.AssetPrefix = DefAssetPrefixName
	c.Ignore = make([]*regexp.Regexp, 0)
	c.Include = make([]*regexp.Regexp, 0)
	return c
}

func (c *Config) validateInput() (err error) {
	uniqPaths := make(map[string]struct{}, len(c.Input))
	newInputs := make([]InputConfig, 0, len(c.Input))

	for _, input := range c.Input {
		input.Path = filepath.Clean(input.Path)
		_, ok := uniqPaths[input.Path]
		if ok {
			continue
		}
		_, err = os.Lstat(input.Path)
		if err != nil {
			return fmt.Errorf("failed to stat input path '%s': %v",
				input.Path, err)
		}
		uniqPaths[input.Path] = struct{}{}
		newInputs = append(newInputs, input)
	}
	if len(newInputs) == 0 {
		return ErrNoInput
	}
	c.Input = newInputs
	return nil
}

// validateOutput will check if output is valid.
//
// (1) If output is empty, set the output directory to,
// (1.1) current working directory if `split` option is used, or
// (1.2) current working directory with default output file output name.
// (2) If output is not empty, check the directory and file write status.
func (c *Config) validateOutput() (err error) {
	// (1)
	if len(c.Output) == 0 {
		if c.Split {
			// (1.1)
			c.Output = c.cwd
		} else {
			// (1.2)
			c.Output = filepath.Join(c.cwd, DefOutputName)
		}

		return nil
	}

	// (2)
	dir, file := filepath.Split(c.Output)

	if dir != "" {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("create output directory: %v", err)
		}
	}

	if len(file) == 0 {
		if !c.Split {
			c.Output = filepath.Join(dir, DefOutputName)
		}
	}

	if c.Split {
		return nil
	}

	var fout *os.File

	fout, err = os.Create(c.Output)
	if err != nil {
		return err
	}

	return fout.Close()
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() (err error) {
	if len(c.cwd) == 0 {
		c.cwd, err = os.Getwd()
		if err != nil {
			return ErrCWD
		}
	}

	if len(c.Package) == 0 {
		return ErrNoPackageName
	}

	err = c.validateInput()
	if err != nil {
		return
	}

	err = c.validateOutput()
	if err != nil {
		return
	}

	return
}
