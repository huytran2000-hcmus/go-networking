package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	t.Logf("bound to %q", listener.Addr())
	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					done <- struct{}{}
					c.Close()
				}()

				data := make([]byte, 1024)
				for {
					n, err := c.Read(data)
					if err != nil {
						if err != io.EOF {
							t.Log(err)
						}
						return
					}
					t.Logf("received: %s\n", data[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}
