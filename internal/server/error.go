package server

import (
	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

// HandlerError simply contains a status code and a message

type HandlerError struct {
	StatusCode  response.StatusCode
	Message     string
	ContentType headers.ContentType
}

func (e *HandlerError) write(w *response.Writer) {

	if w.WriteStatusLine(e.StatusCode) {
		return
	}

	if w.WriteHeaders(headers.GetContentTypeHeaders(len(e.Message), e.ContentType)) {
		return
	}

	if w.WriteBody([]byte(e.Message)) {
		return
	}
}

func newHandlerError(statusCode response.StatusCode, p string) *HandlerError {
	return &HandlerError{StatusCode: statusCode, Message: response.HtmlHandlerErrorMessage(statusCode, p), ContentType: headers.HTML}
}

func internalServerError(err error) *HandlerError {
	return serverError(response.INTERNAL_SERVER_ERROR, err)
}
func badRequestServerError(err error) *HandlerError {
	return serverError(response.BAD_REQUEST, err)
}

func serverError(statusCode response.StatusCode, err error) *HandlerError {
	return &HandlerError{StatusCode: statusCode, Message: err.Error(), ContentType: headers.PLAIN}
}
