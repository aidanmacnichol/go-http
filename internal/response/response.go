package response

import (
	"fmt"
	"go-http/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
	}
	_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprint(contentLength)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for val, key := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", val, key)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
