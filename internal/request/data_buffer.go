package request

import (
	"errors"
	"io"
)

const MINIMALSIZE int = 1

type dataBuffer struct {
	reader      io.Reader
	data        []byte
	readPointer int
	eof         bool
}

func newDataBuffer(reader io.Reader, initialSize int) *dataBuffer {

	if initialSize < MINIMALSIZE {
		initialSize = MINIMALSIZE
	}

	return &dataBuffer{
		reader:      reader,
		data:        make([]byte, initialSize, initialSize),
		readPointer: 0,
		eof:         false,
	}
}

func (d *dataBuffer) current() []byte {
	return d.data[:d.readPointer]
}

func (d *dataBuffer) remove(length int) {
	if length > 0 {
		copy(d.data, d.data[length:])
		d.readPointer -= length
	}
}
func (d *dataBuffer) resizeIfNecessary() {
	if d.readPointer >= len(d.data) {
		newBuf := make([]byte, len(d.data)*2)
		copy(newBuf, d.data)
		d.data = newBuf
	}
}

func (d *dataBuffer) readNext() error {
	n, err := d.reader.Read(d.data[d.readPointer:])
	d.readPointer += n
	return err
}

func (d *dataBuffer) readNextEOF() error {

	d.resizeIfNecessary()

	err := d.readNext()

	if errors.Is(err, io.EOF) {
		d.eof = true
		err = nil
	}

	return err
}
