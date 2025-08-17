package response

import "github.com/Waterbootdev/http-from-tcp/internal/headers"

func (w *Writer) WriteTransferEncodingHeader(trailers []string) bool {

	if w.WriteStatusLine(OK) {
		return true
	}

	if w.setError(writeHeaders(w.writer, headers.GetTransferEncodingHeaders())) {
		return true
	}

	w.status = WriteChunkStatus

	return false
}
