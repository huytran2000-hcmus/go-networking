package udp

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"testing"
	"time"
)

func TestDialUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := echoUDPServer(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.Dial("udp", srv.String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	interupt := []byte("pardon me")
	n, err := interloper.WriteTo(interupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	err = interloper.Close()
	if err != nil {
		t.Fatal(err)
	}

	if l := len(interupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	n, err = client.Write(ping)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(ping); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 256)
	n, err = client.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	got := buf[:n]
	if !bytes.Equal(got, ping) {
		t.Errorf("got %s, want %s", got, ping)
	}

	err = client.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Read(buf)
	if err == nil && errors.Is(err, os.ErrDeadlineExceeded) {
		t.Errorf("expected an error")
	}
}
