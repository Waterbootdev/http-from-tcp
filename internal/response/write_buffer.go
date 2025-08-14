package response

import (
	"bytes"
	"io"
	"log"
)

func WriteBufferOk(w io.Writer, buffer *bytes.Buffer) {
	err := WriteStatusLine(w, OK)

	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	err = WriteDefaultHeaders(w, buffer.Len())

	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	_, err = w.Write(buffer.Bytes())

	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
