package headers

import (
	"strconv"
	"strings"
)

func GetDefaultHeaders(contentLen int) Headers {
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

func GetContentTypeHeaders(contentLen int, contentType ContentType) Headers {
	headers := NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = string(contentType)
	return headers
}

func GetTransferEncodingHeaders() Headers {
	headers := NewHeaders()
	headers["Transfer-Encoding"] = "chunked"
	headers["Content-Type"] = string(PLAIN)
	return headers
}

func GetTransferEncodingTrailerHeaders(trailers []string) Headers {
	headers := GetTransferEncodingHeaders()
	headers["Trailer"] = strings.Join(trailers, ", ")
	return headers
}
