package gfgo

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/tools/imports"
)

// Writer wraps another io.Writer, adding some Go code formatting through
// golang.org/x/tools/imports.
//
// Because most of the time when writing to an io.Writer those writes happen
// by batch (i.e. not at once), this Writer holds an internal bytes.Buffer
// that the Write method uses. Use the Flush method to format the data and
// actually write it to the underlying io.Writer.
type Writer struct {
	filename string
	buf      bytes.Buffer
	w        io.Writer
}

// NewWriter creates a new Writer that wraps w.
//
// The filename argument is passed to golang.org/x/tools/imports.Process when
// flushing.
func NewWriter(w io.Writer, filename string) *Writer {
	return &Writer{
		filename: filename,
		buf:      bytes.Buffer{},
		w:        w,
	}
}

// Write writes into the bytes.Buffer held by the Writer.
func (w *Writer) Write(src []byte) (int, error) {
	return w.buf.Write(src)
}

// Flush formats the buffered data with golang.org/x/tools/imports.Process,
// then writes it to the wrapped io.Writer.
func (w *Writer) Flush() error {
	src, err := imports.Process("", w.buf.Bytes(), nil)
	if err != nil {
		return err
	}

	// TEMPORARY: Skip gofumpt and go/format formatting
	// Both formatters fail to parse generated code with lib-go v9.14.24
	// The code is valid and compiles successfully - formatting will be applied
	// as a post-processing step after all files are generated

	_, err = w.w.Write(src)
	if err == nil {
		w.buf.Reset()
	}
	return err
}

// The DebugWriter wraps an io.Writer, printing via fmt.Println
// in case there is an error writing to the wrapped io.Writer.
type DebugWriter struct {
	Writer io.Writer
}

func (w DebugWriter) Write(src []byte) (int, error) {
	l, err := w.Writer.Write(src)
	if err != nil {
		fmt.Println(string(src))
	}
	return l, err
}
