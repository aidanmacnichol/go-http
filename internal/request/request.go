package request

import (
	"bytes"
	"fmt"
	"io"
)

const bufferSize int = 8

const (
	Initialized int = 0
	Done        int = 1
)

type Request struct {
	RequestLine RequestLine
	ReaderState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

}

func parseRequestLine(request []byte) (*RequestLine, int, error) {

	i := bytes.Index(request, []byte("\r\n"))
	if i == -1 {
		return nil, 0, nil
	}

	startLine := request[:i]
	read := i + len("\r\n")

	splitRequest := bytes.Split(startLine, []byte(" "))
	if len(splitRequest) != 3 {
		return nil, 0, fmt.Errorf("invalid number of header lines in: %s", splitRequest)
	}

	method := splitRequest[0]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, 0, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := splitRequest[1]

	httpVersionParts := bytes.Split(splitRequest[2], []byte("/"))
	if string(httpVersionParts[0]) != "HTTP" {
		return nil, 0, fmt.Errorf("invalid protocol: %s", httpVersionParts[1])
	}
	if string(httpVersionParts[1]) != "1.1" {
		return nil, 0, fmt.Errorf("invalid HTTP version: %s", httpVersionParts[1])
	}

	return &RequestLine{
		Method:        string(method),
		RequestTarget: string(requestTarget),
		HttpVersion:   string(httpVersionParts[1]),
	}, read, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0

outer:
	for {
		switch r.ReaderState {
		case Initialized:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.ReaderState = Done

		case Done:
			break outer
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		ReaderState: Initialized,
	}

	buf := make([]byte, bufferSize)

	readToIndex := 0

	for {
		if request.ReaderState == Done {
			break
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}
	}

	res, err := parseRequestLine(request)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *res,
	}, nil
}
