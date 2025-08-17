package server

import (
	"bytes"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

func handleSucces(w *response.Writer) {
	buffer := bytes.NewBuffer([]byte{})
	_, err := buffer.WriteString(response.HtmlHandlerMessage(response.OK, "Success!", "Your request was an absolute banger."))

	if err != nil {
		return
	}

	if w.WriteBuffer(headers.HTML, buffer) {
		return
	}
}
