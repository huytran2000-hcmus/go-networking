package ch04

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 1 << 20
)

var (
	ErrMaxPayloadSize = errors.New("maximum payload exceeded")
	ErrInvalidType    = errors.New("invalid payload type")
)

type Payload interface {
	String() string
	Byte() []byte
	io.ReaderFrom
	io.WriterTo
}

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, fmt.Errorf("read string payload type: %w", err)
	}

	var payload Payload
	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, ErrInvalidType
	}

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, fmt.Errorf("read a payload from reader: %w", err)
	}

	return payload, nil
}
