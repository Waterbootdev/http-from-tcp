package server

import (
	"errors"

	"github.com/Waterbootdev/http-from-tcp/internal/request"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

const ENDPOINT_PREFIX_HTTPBIN = "/httpbin"
const TARGET_PREFIX_HTTPBIN = "https://httpbin.org"
const DEFAULT_CHUNK_BUFFER_LENGTH = 12 * 1024

type Handler func(w *response.Writer, req *request.Request)

var ExampleHandler Handler = func(w *response.Writer, req *request.Request) {

	taget := req.RequestLine.RequestTarget

	switch taget {
	case "/":
		handleSucces(w)
	case "/video":
		handleVideo(w)
	case "/yourproblem":
		newHandlerError(response.BAD_REQUEST, "Your request honestly kinda sucked.").write(w)
	case "/myproblem":
		newHandlerError(response.INTERNAL_SERVER_ERROR, "Okay, you know what? This one is on me.").write(w)
	default:
		if source := testSwapPrefix(taget, ENDPOINT_PREFIX_HTTPBIN, TARGET_PREFIX_HTTPBIN); source != nil {
			handleTransferEncoding(w, source, DEFAULT_CHUNK_BUFFER_LENGTH)
		} else {
			badRequestServerError(errors.New("invalid prefix")).write(w)
		}
	}
}
