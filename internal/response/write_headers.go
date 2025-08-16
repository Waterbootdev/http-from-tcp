package response

import (
	"io"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
