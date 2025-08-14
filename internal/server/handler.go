package server

import (
	"io"

	"github.com/Waterbootdev/http-from-tcp/internal/request"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

var ExampleHandler Handler = func(w io.Writer, req *request.Request) *HandlerError {

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	}

	w.Write([]byte("All good, frfr\n"))

	return nil
}
