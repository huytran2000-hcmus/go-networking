package udp

import (
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
}
