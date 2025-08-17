package response

import (
	"bytes"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

func (w *Writer) WriteBuffer(contentType headers.ContentType, buffer *bytes.Buffer) bool {

	if w.WriteStatusLine(OK) {
		return true
	}

	if w.WriteHeaders(headers.GetContentTypeHeaders(buffer.Len(), contentType)) {
		return true
	}

	return w.WriteBody(buffer.Bytes())
}
