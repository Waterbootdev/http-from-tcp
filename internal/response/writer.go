package response

import (
	"errors"
	"io"
)

const CONTENT_LENGTH_KEY = "Content-Length"

type WriterStatus int

const (
	WriteStatusLineStatus WriterStatus = iota
	WriteHeadersStatus
	WriteChunkStatus
	WriteBodyStatus
	WriteDoneStatus
	WriteErrorStatus
)

type Writer struct {
	writer io.Writer
	status WriterStatus
	err    []error
}

func NewWriter(ioWriter io.Writer) *Writer {
	return &Writer{writer: ioWriter, status: WriteStatusLineStatus, err: []error{}}
}

func (w *Writer) setError(err error) bool {

	if err == nil {
		return false
	}

	w.err = append(w.err, err)
	w.status = WriteErrorStatus
	return true
}

func (w *Writer) setNError(_ int, err error) bool {
	return w.setError(err)
}

func (w *Writer) setTest(test bool, s string) bool {

	if test {
		w.err = append(w.err, errors.New(s))
		w.status = WriteErrorStatus
	}

	return test
}
