package server

import (
	"fmt"
	"io"
)

// HandlerError simply contains a status code and a message

type HandlerError struct {
	StatusCode int
	Message    string
}

func (e *HandlerError) Write(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "HandlerError: ( %d ) %s\n", e.StatusCode, e.Message)
}
