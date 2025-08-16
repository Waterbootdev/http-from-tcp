package server

import (
	"log"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

// HandlerError simply contains a status code and a message

type HandlerError struct {
	StatusCode  response.StatusCode
	Message     string
	ContentType headers.ContentType
}

func (e *HandlerError) Write(w *response.Writer) {

	err := w.WriteStatusLine(e.StatusCode)

	if err != nil {
		return
	}

	err = w.WriteHeaders(headers.GetContentTypeHeaders(len(e.Message), e.ContentType))

	if err != nil {
		return
	}

	_, err = w.WriteBody([]byte(e.Message))

	if err != nil {
		log.Printf("Error writing handler error response: %v", err)
	}
}

func NewHandlerError(statusCode response.StatusCode, p string) *HandlerError {
	return &HandlerError{StatusCode: statusCode, Message: response.HtmlHandlerErrorMessage(statusCode, p), ContentType: headers.HTML}
}
