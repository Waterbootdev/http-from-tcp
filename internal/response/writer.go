package response

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

const CONTENT_LENGTH_KEY = "Content-Length"

type WriterStatus int

const (
	WriteStatusLineStatus WriterStatus = iota
	WriteHeadersStatus
	WriteBodyStatus
	WriteDoneStatus
)

type Writer struct {
	IoWriter io.Writer
	Status   WriterStatus
}

func NewWriter(ioWriter io.Writer) *Writer {
	return &Writer{IoWriter: ioWriter, Status: WriteStatusLineStatus}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.Status != WriteStatusLineStatus {
		return errors.New("invalid status")
	}

	err := WriteStatusLine(w.IoWriter, statusCode)

	if err != nil {
		return err
	}

	w.Status = WriteHeadersStatus

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.Status != WriteHeadersStatus {
		return errors.New("invalid status")
	}

	err := WriteHeaders(w.IoWriter, headers)

	if err != nil {
		return err
	}

	if headers.IsContentLengthNot(CONTENT_LENGTH_KEY, 0) {
		w.Status = WriteBodyStatus
	} else {
		w.Status = WriteDoneStatus
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.Status != WriteBodyStatus {
		return 0, errors.New("invalid status")
	}

	w.Status = WriteDoneStatus

	return w.IoWriter.Write(p)
}

func (w *Writer) WriteBuffer(contentType ContentType, buffer *bytes.Buffer) error {
	err := w.WriteStatusLine(OK)

	if err != nil {
		return err
	}

	err = w.WriteContentTypeHeaders(buffer.Len(), contentType)

	if err != nil {
		return err
	}

	_, err = w.WriteBody(buffer.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBufferLogError(contentType ContentType, buffer *bytes.Buffer) {

	err := w.WriteBuffer(contentType, buffer)

	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (w *Writer) WriteDefaultHeaders(contentLen int) error {
	return w.WriteHeaders(GetDefaultHeaders(contentLen))
}
func (w *Writer) WriteContentTypeHeaders(contentLen int, contentType ContentType) error {
	return w.WriteHeaders(GetContentTypeHeaders(contentLen, contentType))
}
