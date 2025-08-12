package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Waterbootdev/http-from-tcp/internal/buffer"
	"github.com/Waterbootdev/http-from-tcp/internal/commen"
)

type ParserState int

const (
	Initialized ParserState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
	ParseError  error
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) RequestLineString() string {
	return fmt.Sprintf(`Request line:
- Method: %s
- Target: %s
- Version: %s`, r.Method, r.RequestTarget, r.HttpVersion)
}

func parseRequestLine(line string) (RequestLine, error) {

	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid request line")
	}

	method := parts[0]

	if method != "GET" && method != "POST" {
		return RequestLine{}, errors.New("invalid request method")
	}

	requestTarget := parts[1]

	if !strings.HasPrefix(requestTarget, "/") {
		return RequestLine{}, errors.New("invalid request target")
	}

	protocol := strings.Split(parts[2], "/")

	if len(protocol) != 2 || protocol[0] != "HTTP" {
		return RequestLine{}, errors.New("invalid request protocol")
	}

	httpVersion := protocol[1]

	if httpVersion != "1.1" {
		return RequestLine{}, errors.New("invalid request protocol version")
	}

	return RequestLine{HttpVersion: httpVersion, RequestTarget: requestTarget, Method: method}, nil
}

func (r *Request) parse(data []byte) int {

	crlfIndex := bytes.Index(data, []byte(commen.CRLF))

	if crlfIndex == -1 {
		return 0
	}

	r.ParserState = Done

	r.RequestLine, r.ParseError = parseRequestLine(string(data[:crlfIndex]))

	return crlfIndex + commen.LENGTH_CRLF
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	dataBuffer := buffer.NewDataBuffer(reader, buffer.MINIMALSIZE)

	request := &Request{ParserState: Initialized}

	for request.ParserState != Done {

		err := dataBuffer.ReadNextEOF()

		if err != nil {
			return request, err
		}

		if dataBuffer.EOF {
			request.ParseError = errors.New("unexpected EOF")
			request.ParserState = Done
			break
		}

		dataBuffer.Remove(request.parse(dataBuffer.Current()))
	}

	return request, request.ParseError
}
