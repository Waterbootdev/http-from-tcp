package response

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const MINIMAL_CHUNK_BUFFER_LENGTH = 32
const MAXIMAL_CHUNK_BUFFER_LENGTH = 1024 * 1204

func (w *Writer) writeEndChunk() bool {

	if w.setTest(w.status != WriteChunkStatus, "invalid status") {
		return true
	}

	if w.setNError(w.writer.Write([]byte("0\r\n"))) {
		return true
	}

	w.status = WriteHeadersStatus

	return false
}

func (w *Writer) writeChunk(chunkLength int, chunk []byte) bool {

	if w.setTest(w.status != WriteChunkStatus, "invalid status") {
		return true
	}

	if w.setTest(chunkLength <= 0, "invalid chunk length") {
		return true
	}

	if w.setNError(fmt.Fprintf(w.writer, "%x\r\n", chunkLength)) {
		return true
	}

	if w.setNError(w.writer.Write(chunk)) {
		return true
	}

	if w.setNError(w.writer.Write([]byte("\r\n"))) {
		return true
	}

	return false
}

func (w *Writer) RewriteCunks(reader io.Reader, bufferLength int) *bytes.Buffer {

	readBuffer := bytes.NewBuffer([]byte{})

	chunkBuffer := make([]byte, max(MINIMAL_CHUNK_BUFFER_LENGTH, min(MAXIMAL_CHUNK_BUFFER_LENGTH, bufferLength)))

	for {
		n, err := reader.Read(chunkBuffer)

		if errors.Is(err, io.EOF) {
			break
		}

		if w.setError(err) {
			return nil
		}

		chunk := chunkBuffer[:n]

		if w.setNError(readBuffer.Write(chunk)) {
			return nil
		}

		if w.writeChunk(n, chunk) {
			return nil
		}
	}

	if w.writeEndChunk() {
		return nil
	}

	return readBuffer
}
