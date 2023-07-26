package ch03

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingerAdvancedDeadline(t *testing.T) {
	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	start := time.Now()
	go func() {
		defer func() {
			close(done)
		}()

		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second

		go Pinger(ctx, conn, resetTimer)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 256)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				var nErr net.Error
				ok := errors.As(err, &nErr)
				if !ok || !nErr.Timeout() {
					t.Errorf("want timeout, got %v", err)
				} else {
					t.Log(err)
				}
				return
			}
			t.Logf("receiced: %s (%s)", buf[:n], time.Since(start).Truncate(time.Second))
			resetTimer <- 0
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 256)
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("receiced: %s (%s)", buf[:n], time.Since(start).Truncate(time.Second))
	}

	_, err = conn.Write([]byte("PONG"))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
		}
		t.Logf("receiced: %s (%s)", buf[:n], time.Since(start).Truncate(time.Second))
	}
	<-done
	end := time.Since(start).Truncate(time.Second)
	t.Logf("done (%s)", end)

	want := 9 * time.Second
	if end != want {
		t.Errorf("want EOF at %s, got %s", want, end)
	}
}
