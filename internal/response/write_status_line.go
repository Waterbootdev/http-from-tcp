package response

func (w *Writer) WriteStatusLine(statusCode StatusCode) bool {

	if w.setTest(w.status != WriteStatusLineStatus, "invalid status") {
		return true
	}

	if w.setError(WriteStatusLine(w.writer, statusCode)) {
		return true
	}

	w.status = WriteHeadersStatus

	return false
}
