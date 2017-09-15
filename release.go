// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

// writeOneFileRelease writes the release code file for each file (when splited file).
func writeOneFileRelease(w io.Writer, c *Config, a *Asset) error {
	if err := fileHeader(w); err != nil {
		return err
	}

	if err := writeReleaseAsset(w, c, a); err != nil {
		return err
	}

	return nil
}

// writeRelease writes the release code file for single file.
func writeRelease(w io.Writer, c *Config, toc []Asset) error {
	err := writeReleaseHeader(w, c)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeReleaseAsset(w, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeReleaseHeader writes output file headers.
// This targets release builds.
func writeReleaseHeader(w io.Writer, c *Config) error {
	var err error
	if c.NoCompress {
		if c.NoMemCopy {
			err = headerUncompressedNomemcopy(w)
		} else {
			err = headerUncompressedMemcopy(w)
		}
	} else {
		if c.NoMemCopy {
			err = headerCompressedNomemcopy(w)
		} else {
			err = headerCompressedMemcopy(w)
		}
	}
	if err != nil {
		return err
	}
	return headerReleaseCommon(w)
}

// writeReleaseAsset write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeReleaseAsset(w io.Writer, c *Config, asset *Asset) error {
	fd, err := os.Open(asset.Path)
	if err != nil {
		return err
	}

	defer fd.Close()

	if c.NoCompress {
		if c.NoMemCopy {
			err = uncompressedNomemcopy(w, asset, fd)
		} else {
			err = uncompressedMemcopy(w, asset, fd)
		}
	} else {
		if c.NoMemCopy {
			err = compressedNomemcopy(w, asset, fd)
		} else {
			err = compressedMemcopy(w, asset, fd)
		}
	}
	if err != nil {
		return err
	}
	return assetReleaseCommon(w, c, asset)
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
// Based on https://code.google.com/p/go/source/browse/godoc/static/makestatic.go?repo=tools
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

func fileHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"os"
	"time"
)

`)
	return err
}

func headerCompressedNomemcopy(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

`)
	return err
}

func headerCompressedMemcopy(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

`)
	return err
}

func headerUncompressedNomemcopy(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func bindataRead(data, name string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b, nil
}

`)
	return err
}

func headerUncompressedMemcopy(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
`)
	return err
}

func headerReleaseCommon(w io.Writer) error {
	_, err := fmt.Fprintf(w, `type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

`)
	return err
}

func compressedNomemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, "var _%s = \"\" +\n\t\"", asset.Func)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	gz.Close()

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `"

func %sBytes() ([]byte, error) {
	return bindataRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func compressedMemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, "var _%s = []byte(\n\t\"", asset.Func)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	gz.Close()

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `")

func %sBytes() ([]byte, error) {
	return bindataRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func uncompressedNomemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = "`, asset.Func)
	if err != nil {
		return err
	}

	_, err = io.Copy(&StringWriter{Writer: w}, r)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `"

func %sBytes() ([]byte, error) {
	return bindataRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func uncompressedMemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = []byte(`, asset.Func)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if utf8.Valid(b) && !bytes.Contains(b, []byte{0}) {
		fmt.Fprintf(w, "`%s`", sanitize(b))
	} else {
		fmt.Fprintf(w, "%+q", b)
	}

	_, err = fmt.Fprintf(w, `)

func %sBytes() ([]byte, error) {
	return _%s, nil
}

`, asset.Func, asset.Func)
	return err
}

func assetReleaseCommon(w io.Writer, c *Config, asset *Asset) (err error) {
	fi, err := os.Stat(asset.Path)
	if err != nil {
		return
	}

	mode := uint(fi.Mode())
	modTime := fi.ModTime().Unix()
	size := fi.Size()
	if c.NoMetadata {
		mode = 0
		modTime = 0
		size = 0
	}
	if c.Mode > 0 {
		mode = uint(os.ModePerm) & c.Mode
	}
	if c.ModTime > 0 {
		modTime = c.ModTime
	}

	var md5checksum string
	if c.MD5Checksum {
		var buf []byte

		buf, err = ioutil.ReadFile(asset.Path)
		if err != nil {
			return
		}

		h := md5.New()
		if _, err = h.Write(buf); err != nil {
			return
		}
		md5checksum = fmt.Sprintf("%x", h.Sum(nil))
	}

	_, err = fmt.Fprintf(w, `func %s() (*asset, error) {
	bytes, err := %sBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: %q, size: %d, md5checksum: %q, mode: os.FileMode(%d), modTime: time.Unix(%d, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

`, asset.Func, asset.Func, asset.Name, size, md5checksum, mode, modTime)
	return
}
