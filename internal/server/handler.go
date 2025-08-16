package server

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/request"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

const ENDPOINT_PREFIX_HTTPBIN = "/httpbin"
const TARGET_PREFIX_HTTPBIN = "https://httpbin.org"

type Handler func(w *response.Writer, req *request.Request)

var ExampleHandler Handler = func(w *response.Writer, req *request.Request) {

	taget := req.RequestLine.RequestTarget

	switch taget {
	case "/":
		writeSuccesResponse(w)
	case "/video":
		handleVideo(w, req)
	case "/yourproblem":
		NewHandlerError(response.BAD_REQUEST, "Your request honestly kinda sucked.").Write(w)
	case "/myproblem":
		NewHandlerError(response.INTERNAL_SERVER_ERROR, "Okay, you know what? This one is on me.").Write(w)
	default:
		handlePrefixTargets(w, req, taget)
	}

}

func handleVideo(w *response.Writer, _ *request.Request) {
	video, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		NewHandlerError(response.INTERNAL_SERVER_ERROR, err.Error()).Write(w)
		return
	}
	err = w.WriteStatusLine(response.OK)
	if err != nil {
		log.Printf("Error writing status line: %v", err)
		return
	}
	err = w.WriteHeaders(headers.GetContentTypeHeaders(len(video), headers.MP4))
	if err != nil {
		log.Printf("Error writing headers: %v", err)
		return
	}
	_, err = w.WriteBody(video)
	if err != nil {
		log.Printf("Error writing body: %v", err)
		return
	}
}

func writeSuccesResponse(w *response.Writer) {
	buffer := bytes.NewBuffer([]byte{})
	_, err := buffer.WriteString(response.HtmlHandlerMessage(response.OK, "Success!", "Your request was an absolute banger."))
	if err != nil {
		return
	}
	w.WriteBufferLogError(headers.HTML, buffer)
}

func handlePrefixTargets(w *response.Writer, req *request.Request, taget string) {
	if strings.HasPrefix(taget, ENDPOINT_PREFIX_HTTPBIN) {
		httpbinHandler(w, req, taget)
	}
}

const DEFAULT_CHUNK_BUFFER_LENGTH = 12 * 1024

const XContentSHA256 = "X-Content-SHA256"
const XContentLength = "X-Content-Length"

func httpbinHandler(w *response.Writer, _ *request.Request, taget string) {
	reader, err := httpbinHandlerSendRequest(taget)

	if err != nil {
		NewHandlerError(response.BAD_REQUEST, err.Error()).Write(w)
		return
	}

	defer reader.Close()

	err = w.WriteBeginTransferEncoding([]string{XContentSHA256, XContentLength})

	if err != nil {
		log.Printf("Error writing transfer encoding headers: %v", err)
		return
	}

	body, err := w.RewriteCunks(reader, DEFAULT_CHUNK_BUFFER_LENGTH)

	if err != nil {
		log.Printf("Error rewriting chunks: %v", err)
	}

	header := headers.NewHeaders()

	header[XContentSHA256] = fmt.Sprintf("%x", sha256.Sum256(body.Bytes()))
	header[XContentLength] = strconv.Itoa(body.Len())

	log.Println(header.HeadersString())

	err = w.WriteHeaders(header)

	if err != nil {
		log.Printf("Error writing trailers: %v", err)
	}
}

func httpbinHandlerSendRequest(target string) (io.ReadCloser, error) {

	response, err := http.Get(SwapPrefix(target, ENDPOINT_PREFIX_HTTPBIN, TARGET_PREFIX_HTTPBIN))
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func SwapPrefix(s string, oldPrefix, newPrefix string) string {
	return strings.Replace(s, oldPrefix, newPrefix, 1)
}
