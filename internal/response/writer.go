package response

import (
	"errors"
	"io"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

type WriterStatus int

const (
	WriteStatusLineStatus WriterStatus = iota
	WriteHeadersStatus
	WriteBodyStatus
	WriteDoneStatus
)

type Writer struct {
	ioWriter io.Writer
	Status   WriterStatus
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.Status != WriteStatusLineStatus {
		return errors.New("invalid status")
	}

	err := WriteStatusLine(w.ioWriter, statusCode)

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

	err := WriteHeaders(w.ioWriter, headers)

	if err != nil {
		return err
	}

	if headers.IsContentLengthNot(0) {
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

	return w.ioWriter.Write(p)
}
