package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var caveJS = "../../cave-js"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	fs, err := ioutil.ReadDir(caveJS)
	if err != nil {
		return err
	}
	out, err := os.Create("../../bundle.go")
	if err != nil {
		return err
	}
	_, _ = out.Write([]byte("package cave \n\nconst (\n"))
	for _, f := range fs {
		if !strings.Contains(f.Name(), "bundle.js") {
			continue
		}
		// replace dots with underscores
		_, _ = out.Write([]byte(strings.ReplaceAll(f.Name(), ".", "_") + " = \""))
		f, err := os.Open(filepath.Join(caveJS, f.Name()))
		if err != nil {
			return err
		}
		sw := StringWriter{Writer: out}
		writer := gzip.NewWriter(&sw)
		if _, err = io.Copy(writer, f); err != nil {
			return err
		}
		writer.Close()
		_, _ = out.Write([]byte("\"\n"))
	}
	_, _ = out.Write([]byte(")\n"))
	return nil
}

//https://github.com/go-bindata/go-bindata/blob/master/stringwriter.go
const lowerHex = "0123456789abcdef"

type StringWriter struct {
	io.Writer
	c int
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	buf := []byte(`\x00`)
	var b byte

	for n, b = range p {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		if _, err := w.Writer.Write(buf); err != nil {
			return 0, err
		}
		w.c++
	}
	n++
	return
}
