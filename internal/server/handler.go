package server

import (
	"bytes"

	"github.com/Waterbootdev/http-from-tcp/internal/request"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

var ExampleHandler Handler = func(w *response.Writer, req *request.Request) {

	switch req.RequestLine.RequestTarget {
	case "/":
		writeSuccesResponse(w)
	case "/yourproblem":
		NewHandlerError(response.BAD_REQUEST, "Your request honestly kinda sucked.").Write(w)
	case "/myproblem":
		NewHandlerError(response.INTERNAL_SERVER_ERROR, "Okay, you know what? This one is on me.").Write(w)
	}

}

func writeSuccesResponse(w *response.Writer) {
	buffer := bytes.NewBuffer([]byte{})
	_, err := buffer.WriteString(response.HtmlHandlerMessage(response.OK, "Success!", "Your request was an absolute banger."))
	if err != nil {
		return
	}
	w.WriteBufferLogError(response.HTML, buffer)
}
