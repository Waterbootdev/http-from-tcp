package response

func (w *Writer) WriteBody(p []byte) bool {

	if w.setTest(w.status != WriteBodyStatus, "invalid status") {
		return true
	}

	if w.setNError(w.writer.Write(p)) {
		return true
	}

	w.status = WriteDoneStatus

	return false
}
