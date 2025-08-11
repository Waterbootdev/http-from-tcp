package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF string = "\r\n"

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

func readLines(reader io.Reader) ([]string, error) {

	all, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(all), CRLF), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	lines, err := readLines(reader)

	if err != nil || len(lines) == 0 {
		return nil, err
	}

	requestLine, err := parseRequestLine(lines[0])

	if err != nil {
		return nil, err
	}

	request := Request{RequestLine: requestLine}

	return &request, nil

}
