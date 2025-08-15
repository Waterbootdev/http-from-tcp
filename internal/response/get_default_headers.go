package response

import (
	"io"
	"strconv"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	return GetContentTypeHeaders(contentLen, PLAIN)
}

type ContentType string

const (
	PLAIN      ContentType = "text/plain"
	HTML       ContentType = "text/html"
	CSS        ContentType = "text/css"
	JAVASCRIPT ContentType = "text/javascript"
	JSON       ContentType = "application/json"
)

func GetContentTypeHeaders(contentLen int, contentType ContentType) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = string(contentType)
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func WriteDefaultHeaders(w io.Writer, contentLen int) error {
	return WriteHeaders(w, GetDefaultHeaders(contentLen))
}
func WriteContentTypeHeaders(w io.Writer, contentLen int, contentType ContentType) error {
	return WriteHeaders(w, GetContentTypeHeaders(contentLen, contentType))
}
