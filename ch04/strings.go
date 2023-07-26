package ch04

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type String string

func (m String) String() string {
	return string(m)
}

func (m String) Byte() []byte {
	return []byte(m)
}

func (m String) WriteTo(w io.Writer) (int64, error) {
	var n int64 = 0
	err := binary.Write(w, binary.BigEndian, StringType)
	if err != nil {
		return n, fmt.Errorf("write string payload type: %w", err)
	}
	n += 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, fmt.Errorf("write string payload size: %w", err)
	}
	n += 4

	o, err := w.Write([]byte(m))
	if err != nil {
		return n, fmt.Errorf("write string payload: %w", err)
	}
	n += int64(o)

	return n, err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var n int64 = 0

	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return n, fmt.Errorf("read string payload type: %w", err)
	}
	n += 1

	if typ != StringType {
		return n, ErrInvalidType
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, fmt.Errorf("read string payload size: %w", err)
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	buf := make([]byte, size)
	o, err := r.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, fmt.Errorf("read string payload: %w", err)
	}
	n += int64(o)
	*m = String(buf)

	return n, err
}
