package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Waterbootdev/http-from-tcp/internal/buffer"
	"github.com/Waterbootdev/http-from-tcp/internal/commen"
	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

type ParserState int

const (
	Initialized ParserState = iota
	RequestStateParsingHeaders
	RequestStateParsingBody
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	ParserState ParserState
	ParseError  error
}

func (r *Request) BodyString() string {

	var buffer bytes.Buffer
	buffer.WriteString("Body:\r\n")
	buffer.Write(r.Body)
	buffer.WriteString("\r\n")
	return buffer.String()
}

func (r *Request) checkValidEOF() {

	if r.ParserState != RequestStateParsingBody {
		r.ParseError = errors.New("unexpected EOF")
		return
	}

	if r.Headers.IsContentLengthNot(len(r.Body)) {
		r.ParseError = errors.New("content length does not match body length")
		return
	}
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

func splitRequestLine(line string) (string, string, string, error) {

	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return "", "", "", errors.New("invalid request line")
	}

	method := parts[0]

	if method != "GET" && method != "POST" {
		return "", "", "", errors.New("invalid request method")
	}

	requestTarget := parts[1]

	if !strings.HasPrefix(requestTarget, "/") {
		return "", "", "", errors.New("invalid request target")
	}

	protocol := strings.Split(parts[2], "/")

	if len(protocol) != 2 || protocol[0] != "HTTP" {
		return "", "", "", errors.New("invalid request protocol")
	}

	httpVersion := protocol[1]

	if httpVersion != "1.1" {
		return "", "", "", errors.New("invalid request protocol version")
	}

	return httpVersion, requestTarget, method, nil
}

func (r *RequestLine) parseRequestLine(data []byte) (int, bool, error) {

	crlfIndex := bytes.Index(data, []byte(commen.CRLF))

	if crlfIndex == -1 {
		return 0, false, nil
	}

	var err error

	r.HttpVersion, r.RequestTarget, r.Method, err = splitRequestLine(string(data[:crlfIndex]))

	if err != nil {
		return 0, false, err
	}

	return crlfIndex + commen.LENGTH_CRLF, true, nil
}

func (r *Request) parse(data []byte) int {
	var n int = 0
	var done bool = false

	switch r.ParserState {

	case Initialized:
		n, done, r.ParseError = r.RequestLine.parseRequestLine(data)
	case RequestStateParsingHeaders:
		n, done, r.ParseError = r.Headers.Parse(data)
	case RequestStateParsingBody:
		r.Body = append(r.Body, data...)
		n = len(data)
	}

	if r.ParseError != nil {
		r.ParserState = Done
		return 0
	}

	if done {
		r.NextState()
	}

	return n
}

func (r *Request) NextState() {

	switch r.ParserState {

	case Initialized:
		r.ParserState = RequestStateParsingHeaders
	case RequestStateParsingHeaders:
		r.ParserState = RequestStateParsingBody
		r.ParserState = Done
	case RequestStateParsingBody:
		r.ParserState = Done
	}

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	dataBuffer := buffer.NewDataBuffer(reader, buffer.MINIMALSIZE)

	request := &Request{ParserState: Initialized, Headers: headers.NewHeaders(), Body: []byte{}}

	for request.ParserState != Done {

		err := dataBuffer.ReadNextEOF()

		if err != nil {
			return request, err
		}

		if dataBuffer.EOF {
			request.checkValidEOF()
			request.ParserState = Done
			break
		}

		dataBuffer.Remove(request.parse(dataBuffer.Current()))

	}

	return request, request.ParseError
}

func (r *Request) Print() {
	fmt.Println(r.RequestLine.RequestLineString())
	fmt.Println()
	fmt.Println(r.Headers.HeadersString())
	fmt.Println()
	fmt.Println(r.BodyString())
	fmt.Println()
}
