package response

import (
	"io"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

func writeHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) bool {

	if w.setTest(w.status != WriteHeadersStatus, "invalid status") {
		return true
	}

	if w.setError(writeHeaders(w.writer, headers)) {
		return true
	}

	if headers.IsContentLengthNot(CONTENT_LENGTH_KEY, 0) {
		w.status = WriteBodyStatus
	} else {
		w.status = WriteDoneStatus
	}

	return false
}
