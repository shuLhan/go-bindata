// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shuLhan/go-bindata"
)

var (
	argIgnore  []string
	argVersion bool
	argPrefix  string
	cfg        *bindata.Config
)

func main() {
	initArgs()

	parseArgs()

	err := bindata.Translate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata: %v\n", err)
		os.Exit(1)
	}
}

//
// initArgs will initialize all command line arguments.
//
func initArgs() {
	cfg = bindata.NewConfig()

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input directories>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&argVersion, "version", false, "Displays version information.")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.BoolVar(&cfg.Dev, "dev", cfg.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.")
	flag.BoolVar(&cfg.MD5Checksum, "md5checksum", cfg.MD5Checksum, "MD5 checksums will be calculated for assets.")
	flag.BoolVar(&cfg.NoCompress, "nocompress", cfg.NoCompress, "Assets will *not* be GZIP compressed when this flag is specified.")
	flag.BoolVar(&cfg.NoMemCopy, "nomemcopy", cfg.NoMemCopy, "Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	flag.BoolVar(&cfg.NoMetadata, "nometadata", cfg.NoMetadata, "Assets will not preserve size, mode, and modtime info.")
	flag.Int64Var(&cfg.ModTime, "modtime", cfg.ModTime, "Optional modification unix timestamp override for all files.")
	flag.StringVar(&argPrefix, "prefix", "", "Optional path prefix to strip off asset names.")
	flag.StringVar(&cfg.Output, "o", cfg.Output, "Optional name of the output file to be generated.")
	flag.StringVar(&cfg.Package, "pkg", cfg.Package, "Package name to use in the generated code.")
	flag.StringVar(&cfg.Tags, "tags", cfg.Tags, "Optional set of build tags to include.")
	flag.UintVar(&cfg.Mode, "mode", cfg.Mode, "Optional file mode override for all files.")
	flag.Var((*AppendSliceValue)(&argIgnore), "ignore", "Regex pattern to ignore")
}

//
// parseArgs creates a new, filled configuration instance by reading and parsing
// command line options.
//
// This function exits the program with an error, if any of the command line
// options are incorrect.
//
func parseArgs() {
	flag.Parse()

	if argPrefix != "" {
		var err error
		cfg.Prefix, err = regexp.Compile(argPrefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to understand -prefix regex pattern.\n")
			os.Exit(1)
		}
	} else {
		cfg.Prefix = nil
	}

	patterns := make([]*regexp.Regexp, 0)
	for _, pattern := range argIgnore {
		patterns = append(patterns, regexp.MustCompile(pattern))
	}
	cfg.Ignore = patterns

	if argVersion {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	// Make sure we have input paths.
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Missing <input dir>\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create input configurations.
	cfg.Input = make([]bindata.InputConfig, flag.NArg())
	for i := range cfg.Input {
		cfg.Input[i] = parseInput(flag.Arg(i))
	}

	// Change pkg to containing directory of output. If output flag is set and package flag is not.
	pkgSet := false
	outputSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "pkg" {
			pkgSet = true
		}
		if f.Name == "o" {
			outputSet = true
		}
	})
	if outputSet && !pkgSet {
		pkg := filepath.Base(filepath.Dir(cfg.Output))
		if pkg != "." && pkg != "/" {
			cfg.Package = pkg
		}
	}
}

//
// parseInput determines whether the given path has a recrusive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
//
func parseInput(path string) bindata.InputConfig {
	if strings.HasSuffix(path, "/...") {
		return bindata.InputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	}

	return bindata.InputConfig{
		Path: filepath.Clean(path),
	}
}
