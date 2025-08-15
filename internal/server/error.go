package server

import (
	"fmt"
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

func htmlBody(h1 string, p string) string {
	return fmt.Sprintf(`<body>
    <h1>%s</h1>
	<p>%s</p>
  </body>`, h1, p)

}

func HtmlHandlerErrorMessage(statusCode response.StatusCode, h1 string, p string) string {
	return fmt.Sprintf(`<html>
  %s
  %s
</html>`, response.HtmlHead(statusCode), htmlBody(h1, p))
}
