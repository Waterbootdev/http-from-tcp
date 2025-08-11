package request

import (
	"errors"
	"io"
)

const MINIMALSIZE int = 1

type dataBuffer struct {
	reader      io.Reader
	data        []byte
	numberBytes int
	eof         bool
}

func newDataBuffer(reader io.Reader, initialSize int) *dataBuffer {

	if initialSize < MINIMALSIZE {
		initialSize = MINIMALSIZE
	}

	return &dataBuffer{
		reader:      reader,
		data:        make([]byte, initialSize, initialSize),
		numberBytes: 0,
		eof:         false,
	}
}

func (d *dataBuffer) current() []byte {
	return d.data[:d.numberBytes]
}

func (d *dataBuffer) remove(numberBytesToRemove int) {
	if numberBytesToRemove > 0 {
		copy(d.data, d.data[numberBytesToRemove:])
		d.numberBytes -= numberBytesToRemove
	}
}
func (d *dataBuffer) resizeIfNecessary() {
	if d.numberBytes >= len(d.data) {
		newBuf := make([]byte, len(d.data)*2)
		copy(newBuf, d.data)
		d.data = newBuf
	}
}

func (d *dataBuffer) readNext() error {
	numberBytesReaded, err := d.reader.Read(d.data[d.numberBytes:])
	d.numberBytes += numberBytesReaded
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
