package udp

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestEchoUDPServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverAddr, err := echoUDPServer(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("cho khai")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 256)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != serverAddr.String() {
		t.Fatalf("got packet from %s, want from %s", addr.String(), serverAddr.String())
	}

	got := buf[:n]
	t.Logf("%s", got)
	if !bytes.Equal(got, msg) {
		t.Errorf("got %s, want %s", got, msg)
	}
}
