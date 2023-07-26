package ch04

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Binary []byte

func (m Binary) String() string {
	return string(m)
}

func (m Binary) Byte() []byte {
	return m
}

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	var n int64 = 0
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return n, fmt.Errorf("write binary payload type: %w", err)
	}
	n += 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, fmt.Errorf("write binary payload size: %w", err)
	}
	n += 4

	o, err := w.Write(m)
	if err != nil {
		return n, fmt.Errorf("write binary payload: %w", err)
	}
	n += int64(o)

	return n, nil
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var n int64 = 0

	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return n, fmt.Errorf("read binary payload type: %w", err)
	}
	n += 1

	if typ != BinaryType {
		return n, ErrInvalidType
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, fmt.Errorf("read binary payload size: %w", err)
	}
	n += 4

	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size)
	o, err := r.Read(*m)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, fmt.Errorf("read binary payload: %w", err)
	}
	n += int64(o)

	return n, err
}
