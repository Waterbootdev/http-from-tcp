package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

const XContentSHA256 = "X-Content-SHA256"
const XContentLength = "X-Content-Length"

func getSourceReader(w *response.Writer, source *string) io.ReadCloser {
	resp, err := http.Get(*source)
	if err != nil {
		internalServerError(err).write(w)
		return nil
	}
	return resp.Body
}

func handleTransferEncoding(w *response.Writer, source *string, bufferLength int) {

	reader := getSourceReader(w, source)

	if reader == nil {
		return
	}

	defer reader.Close()

	if w.WriteTransferEncodingHeader([]string{XContentSHA256, XContentLength}) {
		return
	}

	body := w.RewriteCunks(reader, bufferLength)

	if body == nil {
		return
	}

	header := headers.NewHeaders()

	header[XContentSHA256] = fmt.Sprintf("%x", sha256.Sum256(body.Bytes()))
	header[XContentLength] = strconv.Itoa(body.Len())

	if w.WriteHeaders(header) {
		return
	}
}
