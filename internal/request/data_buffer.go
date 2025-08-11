package request

import (
	"errors"
	"io"
)

type dataBuffer struct {
	reader      io.Reader
	data        []byte
	readPointer int
	eof         bool
}

func (d *dataBuffer) init(reader io.Reader) {
	d.reader = reader
	d.data = make([]byte, BYTES_PER_CHUNK, BYTES_PER_CHUNK)
	d.readPointer = 0
	d.eof = false
}
func (d *dataBuffer) current() []byte {
	return d.data[:d.readPointer]
}

func (d *dataBuffer) remove(length int) {
	copy(d.data, d.data[length:])
	d.readPointer -= length
}
func (d *dataBuffer) resizeIfNecessary() {
	if d.readPointer >= len(d.data) {
		newBuf := make([]byte, len(d.data)*2)
		copy(newBuf, d.data)
		d.data = newBuf
	}
}

func (d *dataBuffer) readChunk() error {

	d.resizeIfNecessary()

	n, err := d.reader.Read(d.data[d.readPointer:])

	d.readPointer += n

	if err != nil {
		if errors.Is(err, io.EOF) {
			d.eof = true
			return nil
		}

		return err
	}

	return nil
}
