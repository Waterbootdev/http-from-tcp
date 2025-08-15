package response

import (
	"bytes"
	"io"
	"log"
)

func WriteBuffer(w io.Writer, buffer *bytes.Buffer) error {
	err := WriteStatusLine(w, OK)

	if err != nil {
		return err
	}

	err = WriteDefaultHeaders(w, buffer.Len())

	if err != nil {
		return err
	}

	_, err = w.Write(buffer.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func WriteBufferLogError(w io.Writer, buffer *bytes.Buffer) {

	err := WriteBuffer(w, buffer)

	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
