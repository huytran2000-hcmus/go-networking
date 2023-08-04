package udp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
)

func echoUDPServer(ctx context.Context, addr string) (net.Addr, error) {
	server, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen for packet at %s: %w", addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			server.Close()
		}()

		buf := make([]byte, 256)
		for {
			n, clientAddr, err := server.ReadFrom(buf)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("read from udp connection: %v", err)
				}
				return
			}

			_, err = server.WriteTo(buf[:n], clientAddr)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("write from udp connection: %v", err)
				}
				return
			}
		}
	}()

	return server.LocalAddr(), nil
}
