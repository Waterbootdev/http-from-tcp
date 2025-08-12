package buffer

import (
	"errors"
	"io"
)

const MINIMALSIZE int = 1

type DataBuffer struct {
	reader      io.Reader
	data        []byte
	NumberBytes int
	EOF         bool
}

func NewDataBuffer(reader io.Reader, initialSize int) *DataBuffer {

	if initialSize < MINIMALSIZE {
		initialSize = MINIMALSIZE
	}

	return &DataBuffer{
		reader:      reader,
		data:        make([]byte, initialSize, initialSize),
		NumberBytes: 0,
		EOF:         false,
	}
}

func (d *DataBuffer) Current() []byte {
	return d.data[:d.NumberBytes]
}

func (d *DataBuffer) Remove(numberBytesToRemove int) {
	if numberBytesToRemove > 0 {
		copy(d.data, d.data[numberBytesToRemove:])
		d.NumberBytes -= numberBytesToRemove
	}
}
func (d *DataBuffer) resizeIfNecessary() {
	if d.NumberBytes >= len(d.data) {
		newBuf := make([]byte, len(d.data)*2)
		copy(newBuf, d.data)
		d.data = newBuf
	}
}

func (d *DataBuffer) readNext() error {
	numberBytesReaded, err := d.reader.Read(d.data[d.NumberBytes:])
	d.NumberBytes += numberBytesReaded
	return err
}

func (d *DataBuffer) ReadNextEOF() error {

	d.resizeIfNecessary()

	err := d.readNext()

	if errors.Is(err, io.EOF) {
		d.EOF = true
		err = nil
	}

	return err
}
