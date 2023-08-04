package udp

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestListenPackageUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := echoUDPServer(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.ListenPacket("udp", "127.0.0.1:")
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

	buf := make([]byte, 256)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != interloper.LocalAddr().String() {
		t.Fatalf("got reply from %s, want from %s", addr, srv.String())
	}

	got := buf[:n]
	if !bytes.Equal(got, interupt) {
		t.Errorf("got %s, want %s", got, interupt)
	}

	ping := []byte("ping")
	_, err = client.WriteTo(ping, srv)
	if err != nil {
		t.Fatal(err)
	}

	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != srv.String() {
		t.Fatalf("got reply from %s, want from %s", addr, srv.String())
	}

	got = buf[:n]
	if !bytes.Equal(got, ping) {
		t.Errorf("got %s, want %s", got, ping)
	}
}
