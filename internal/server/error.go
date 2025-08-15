package server

import (
	"io"
	"log"

	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

// HandlerError simply contains a status code and a message

type HandlerError struct {
	StatusCode  response.StatusCode
	Message     string
	ContentType response.ContentType
}

func (e *HandlerError) Write(w io.Writer) {

	err := response.WriteStatusLine(w, e.StatusCode)

	if err != nil {
		return
	}

	err = response.WriteHeaders(w, response.GetContentTypeHeaders(len(e.Message), e.ContentType))

	if err != nil {
		return
	}

	_, err = w.Write([]byte(e.Message))

	if err != nil {
		log.Printf("Error writing handler error response: %v", err)
	}
}
