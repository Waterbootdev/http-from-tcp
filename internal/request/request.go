package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
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

const CRLF string = "\r\n"
const LENGTH_CRLF int = len(CRLF)

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

	crlfIndex := bytes.Index(data, []byte(CRLF))

	if crlfIndex == -1 {
		return 0
	}

	r.ParserState = Done

	r.RequestLine, r.ParseError = parseRequestLine(string(data[:crlfIndex]))

	return crlfIndex + LENGTH_CRLF
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	dataBuffer := newDataBuffer(reader, MINIMALSIZE)

	request := &Request{ParserState: Initialized}

	for request.ParserState != Done {

		err := dataBuffer.readNextEOF()

		if err != nil {
			return request, err
		}

		if dataBuffer.eof {
			request.ParserState = Done
			break
		}

		dataBuffer.remove(request.parse(dataBuffer.current()))
	}

	return request, request.ParseError
}
