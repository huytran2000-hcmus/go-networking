package ch03

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		err = conn.SetDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		var nErr net.Error
		ok := errors.As(err, &nErr)
		if !ok || !nErr.Timeout() {
			t.Errorf("want timeout error, got: %v", err)
		}

		done <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-done

	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 10)
	_, err = conn.Read(buf)
	if err != io.EOF {
		t.Errorf("want EOF error, got %v", err)
	}
}
