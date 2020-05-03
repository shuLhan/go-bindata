// Copyright 2018 The go-bindata Authors. All rights reserved.
// Use of this source code is governed by a CC0 1.0 Universal (CC0 1.0)
// Public Domain Dedication license that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"compress/gzip"
	"crypto/md5" //nolint: gas
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

// writeOneFileRelease writes the release code file for each file (when splited file).
func writeOneFileRelease(w io.Writer, c *Config, ast *asset) (err error) {
	_, err = fmt.Fprint(w, tmplImport)
	if err != nil {
		return
	}

	return writeReleaseAsset(w, c, ast)
}

// writeRelease writes the release code file for single file.
func writeRelease(w io.Writer, c *Config, keys []string, toc map[string]*asset) (err error) {
	err = writeReleaseHeader(w, c)
	if err != nil {
		return err
	}

	for _, key := range keys {
		ast := toc[key]
		err = writeReleaseAsset(w, c, ast)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeReleaseHeader writes output file headers.
// This targets release builds.
func writeReleaseHeader(w io.Writer, c *Config) (err error) {
	if c.NoCompress {
		if c.NoMemCopy {
			_, err = fmt.Fprint(w, tmplImportNocompressNomemcopy)
		} else {
			_, err = fmt.Fprint(w, tmplImportNocompressMemcopy)
		}
	} else {
		if c.NoMemCopy {
			_, err = fmt.Fprint(w, tmplImportCompressNomemcopy)
		} else {
			_, err = fmt.Fprint(w, tmplImportCompressMemcopy)
		}
	}
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, tmplReleaseHeader)

	return err
}

// writeReleaseAsset write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeReleaseAsset(w io.Writer, c *Config, ast *asset) (err error) {
	fd, err := os.Open(ast.path)
	if err != nil {
		return
	}

	if c.NoCompress {
		if c.NoMemCopy {
			err = nocompressNomemcopy(w, ast, fd)
		} else {
			err = nocompressMemcopy(w, ast, fd)
		}
	} else {
		if c.NoMemCopy {
			err = compressNomemcopy(w, ast, fd)
		} else {
			err = compressMemcopy(w, ast, fd)
		}
	}
	if err != nil {
		_ = fd.Close()
		return
	}

	err = fd.Close()
	if err != nil {
		return
	}

	return assetReleaseCommon(w, c, ast)
}

//nolint: gochecknoglobals
var (
	backquote = []byte("`")
	bom       = []byte("\xEF\xBB\xBF")
)

// sanitize prepares a valid UTF-8 string as a raw string constant.
// Based on https://code.google.com/p/go/source/browse/godoc/static/makestatic.go?repo=tools
func sanitize(b []byte) []byte {
	var chunks [][]byte
	for i, b := range bytes.Split(b, backquote) {
		if i > 0 {
			chunks = append(chunks, backquote)
		}
		for j, c := range bytes.Split(b, bom) {
			if j > 0 {
				chunks = append(chunks, bom)
			}
			if len(c) > 0 {
				chunks = append(chunks, c)
			}
		}
	}

	var buf bytes.Buffer
	sanitizeChunks(&buf, chunks)
	return buf.Bytes()
}

func sanitizeChunks(buf *bytes.Buffer, chunks [][]byte) {
	n := len(chunks)
	if n >= 2 {
		buf.WriteString("(")
		sanitizeChunks(buf, chunks[:n/2])
		buf.WriteString(" + ")
		sanitizeChunks(buf, chunks[n/2:])
		buf.WriteString(")")
		return
	}
	b := chunks[0]
	if bytes.Equal(b, backquote) {
		buf.WriteString("\"`\"")
		return
	}
	if bytes.Equal(b, bom) {
		buf.WriteString(`"\xEF\xBB\xBF"`)
		return
	}
	buf.WriteString("`")
	buf.Write(b)
	buf.WriteString("`")
}

func compressNomemcopy(w io.Writer, ast *asset, r io.Reader) (err error) {
	_, err = fmt.Fprintf(w, `var _%s = "`, ast.funcName)
	if err != nil {
		return
	}

	gz := gzip.NewWriter(&stringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	if err != nil {
		_ = gz.Close()
		return
	}

	err = gz.Close()
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, tmplFuncCompressNomemcopy, ast.funcName,
		ast.funcName, ast.name)

	return
}

func compressMemcopy(w io.Writer, ast *asset, r io.Reader) (err error) {
	_, err = fmt.Fprintf(w, `var _%s = []byte("`, ast.funcName)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&stringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	if err != nil {
		_ = gz.Close()
		return err
	}

	err = gz.Close()
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, tmplFuncCompressMemcopy, ast.funcName,
		ast.funcName, ast.name)

	return
}

func nocompressNomemcopy(w io.Writer, ast *asset, r io.Reader) (err error) {
	_, err = fmt.Fprintf(w, `var _%s = "`, ast.funcName)
	if err != nil {
		return
	}

	_, err = io.Copy(&stringWriter{Writer: w}, r)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, tmplFuncNocompressNomemcopy, ast.funcName,
		ast.funcName, ast.name)

	return
}

func nocompressMemcopy(w io.Writer, ast *asset, r io.Reader) (err error) {
	_, err = fmt.Fprintf(w, `var _%s = []byte(`, ast.funcName)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	if utf8.Valid(b) && !bytes.Contains(b, []byte{0}) {
		_, err = w.Write(sanitize(b))
	} else {
		_, err = fmt.Fprintf(w, "%+q", b)
	}
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, tmplFuncNocompressMemcopy, ast.funcName,
		ast.funcName)

	return
}

// nolint: gas
func assetReleaseCommon(w io.Writer, c *Config, ast *asset) (err error) {
	fi, err := os.Stat(ast.path)
	if err != nil {
		return err
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

		buf, err = ioutil.ReadFile(ast.path)
		if err != nil {
			return err
		}

		h := md5.New()
		if _, err = h.Write(buf); err != nil {
			return err
		}
		md5checksum = fmt.Sprintf("%x", h.Sum(nil))
	}

	_, err = fmt.Fprintf(w, tmplReleaseCommon, ast.funcName, ast.funcName,
		ast.name, size, md5checksum, mode, modTime)

	return err
}
