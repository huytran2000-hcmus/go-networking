package ch04

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	b1 := Binary("Clear is better than clever")
	b2 := Binary("Don't panic")
	s1 := String("Errors are values")
	payloads := []Payload{&b1, &b2, &s1}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for _, p := range payloads {
			_, err := p.WriteTo(conn)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
		return
	}
	defer conn.Close()

	for _, want := range payloads {
		got, err := decode(conn)
		if err != nil {
			t.Fatal(err)
			return
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}

		t.Logf("[%T] %[1]q", got)
	}
}

func TestMaxPayloadSize(t *testing.T) {
	var buf bytes.Buffer
	err := buf.WriteByte(BinaryType)
	if err != nil {
		t.Fatal(err)
	}

	err = binary.Write(&buf, binary.BigEndian, uint32(1<<10))
	if err != nil {
		t.Fatal(err)
	}

	var b Binary
	_, err = b.ReadFrom(&buf)
	if err != ErrMaxPayloadSize {
		t.Errorf("got error %v, want %v", err, ErrMaxPayloadSize)
	}
}
