package response

import (
	"fmt"
	"go-http/internal/headers"
)

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLength))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
