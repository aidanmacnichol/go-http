package request

import (
	"bytes"
	"fmt"
	"io"
	"errors"
	"strings"
)

const bufferSize int = 8
const crlf = "\r\n"

type requestState int

const (
	readerStateInitialized requestState = iota 
	readerStateDone        
)

type Request struct {
	RequestLine RequestLine
	state requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state: readerStateInitialized,
	}

	buf := make([]byte, bufferSize, bufferSize)
	readToIdx := 0

	for req.state != readerStateDone {
		if readToIdx >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = readerStateDone
				break
			}
			return nil, err
		}
		readToIdx += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIdx])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIdx -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil

}

func requestLineFromString(str string) (*RequestLine, error) {
	splitRequest := strings.Split(str, " ")
	if len(splitRequest) != 3 {
		return nil, fmt.Errorf("invalid number of header lines in: %s", splitRequest)
	}

	method := splitRequest[0]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := splitRequest[1]

	httpVersionParts := strings.Split(splitRequest[2], "/")
	if string(httpVersionParts[0]) != "HTTP" {
		return nil, fmt.Errorf("invalid protocol: %s", httpVersionParts[1])
	}
	if string(httpVersionParts[1]) != "1.1" {
		return nil, fmt.Errorf("invalid HTTP version: %s", httpVersionParts[1])
	}

	return &RequestLine{
		Method:        string(method),
		RequestTarget: string(requestTarget),
		HttpVersion:   string(httpVersionParts[1]),
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case readerStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = readerStateDone
		return n, nil
	case readerStateDone:
		return 0, fmt.Errorf("error: trying to reqad data in a done state")
	default:
		return 0, fmt.Errorf("unkown state")
	}
}

