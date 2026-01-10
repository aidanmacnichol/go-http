package response

import (
	"fmt"
	"go-http/internal/headers"
	"io"
)

type WriterState int

const (
	writerStateHeader WriterState = iota
	writerStateStatusLine
	writerStateBody
)

type Writer struct {
	w     io.Writer
	state WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     w,
		state: writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("invalid state for writing status line: %d", w.state)
	}

	defer func() { w.state = writerStateHeader }()
	_, err := w.w.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateHeader {
		return fmt.Errorf("invalid state for writing headers: %d", w.state)
	}

	defer func() { w.state = writerStateBody }()

	for key, val := range h {
		_, err := w.w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("invalid state for writing body: %d", w.state)
	}
	return w.w.Write(p)
}
