package response

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Waterbootdev/http-from-tcp/internal/headers"
)

const CONTENT_LENGTH_KEY = "Content-Length"

type WriterStatus int

const (
	WriteStatusLineStatus WriterStatus = iota
	WriteHeadersStatus
	WriteChunkStatus
	WriteBodyStatus
	WriteDoneStatus
)

type Writer struct {
	IoWriter io.Writer
	Status   WriterStatus
}

const MINIMAL_CHUNK_BUFFER_LENGTH = 32
const MAXIMAL_CHUNK_BUFFER_LENGTH = 1024 * 1204

func (w *Writer) WriteTrailers(trailers []string) error {
	return nil
}

func (w *Writer) RewriteCunks(reader io.Reader, bufferLength int) (readBuffer *bytes.Buffer, err error) {

	readBuffer = bytes.NewBuffer([]byte{})

	buffer := make([]byte, max(MINIMAL_CHUNK_BUFFER_LENGTH, min(MAXIMAL_CHUNK_BUFFER_LENGTH, bufferLength)))

	for {
		n, err := reader.Read(buffer)

		if err != nil {

			if errors.Is(err, io.EOF) {
				_, err = w.WriteEndChunk()

				if err != nil {
					return readBuffer, err
				}

				break
			}

			w.Status = WriteDoneStatus

			return readBuffer, err
		}

		readBuffer.Write(buffer[:n])

		_, err = w.WriteChunk(n, buffer[:n])

		if err != nil {
			return readBuffer, err
		}
	}

	return readBuffer, nil
}

func (w *Writer) WriteEndChunk() (int, error) {
	if w.Status != WriteChunkStatus {
		return 0, errors.New("invalid status")
	}

	n, err := w.IoWriter.Write([]byte("0\r\n"))

	if err != nil {
		w.Status = WriteDoneStatus
		return n, err
	}

	w.Status = WriteHeadersStatus

	return n, nil
}

func (w *Writer) WriteChunk(chunkLength int, chunk []byte) (int, error) {
	if w.Status != WriteChunkStatus {
		return 0, errors.New("invalid status")
	}

	if chunkLength <= 0 {
		return 0, errors.New("invalid chunk length")
	}

	tn := 0

	n, err := fmt.Fprintf(w.IoWriter, "%x\r\n", chunkLength)
	tn += n

	if err != nil {
		return tn, err
	}

	n, err = w.IoWriter.Write(chunk)
	tn += n

	if err != nil {
		return tn, err
	}

	n, err = w.IoWriter.Write([]byte("\r\n"))
	tn += n

	if err != nil {
		return tn, err
	}

	return tn, nil
}

func NewWriter(ioWriter io.Writer) *Writer {
	return &Writer{IoWriter: ioWriter, Status: WriteStatusLineStatus}
}

func (w *Writer) WriteBeginTransferEncoding(trailers []string) error {
	err := w.WriteStatusLine(OK)

	if err != nil {
		return err
	}

	err = WriteHeaders(w.IoWriter, headers.GetTransferEncodingTrailerHeaders(trailers))

	if err != nil {
		return err
	}

	w.Status = WriteChunkStatus
	return nil
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

func (w *Writer) WriteBuffer(contentType headers.ContentType, buffer *bytes.Buffer) error {
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

func (w *Writer) WriteBufferLogError(contentType headers.ContentType, buffer *bytes.Buffer) {

	err := w.WriteBuffer(contentType, buffer)

	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (w *Writer) WriteDefaultHeaders(contentLen int) error {
	return w.WriteHeaders(headers.GetDefaultHeaders(contentLen))
}
func (w *Writer) WriteContentTypeHeaders(contentLen int, contentType headers.ContentType) error {
	return w.WriteHeaders(headers.GetContentTypeHeaders(contentLen, contentType))
}
