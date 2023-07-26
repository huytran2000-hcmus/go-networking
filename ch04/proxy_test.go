package ch04

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
)

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup
	server := setupTestServer(t, &wg)
	proxyServer := setupTestProxyServer(t, &wg, server)
	client, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct {
		message, want string
	}{
		{"PING", "PONG"},
		{"pong", "pong"},
		{"echo", "echo"},
	}
	for _, tt := range msgs {
		t.Run(fmt.Sprintf("%s->%s", tt.message, tt.want), func(t *testing.T) {
			_, err = client.Write([]byte(tt.message))
			if err != nil {
				t.Error(err)
				return
			}

			buf := make([]byte, 256)
			n, err := client.Read(buf)
			if err != nil {
				t.Error(err)
				return
			}

			got := string(buf[:n])
			t.Logf("%q -> proxy -> %q", tt.message, got)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}

	client.Close()
	proxyServer.Close()
	server.Close()
	wg.Wait()
}

func setupTestProxyServer(t *testing.T, wg *sync.WaitGroup, server net.Listener) net.Listener {
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		for {
			inConn, err := proxyServer.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					t.Error(err)
				}
				return
			}
			defer inConn.Close()

			go func() {
				outConn, err := net.Dial("tcp", server.Addr().String())
				if err != nil {
					t.Error(err)
					return
				}
				defer outConn.Close()

				err = proxy(inConn, outConn)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}()
		}
	}()

	return proxyServer
}

func setupTestServer(t *testing.T, wg *sync.WaitGroup) net.Listener {
	wg.Add(1)
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					t.Error(err)
				}
				return
			}

			go func(conn net.Conn) {
				defer conn.Close()

				for {
					buf := make([]byte, 256)
					n, err := conn.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}

						return
					}

					msg := string(buf[:n])
					switch msg {
					case "PING":
						_, err = conn.Write([]byte("PONG"))
					default:
						_, err = conn.Write(buf[:n])
					}
					if err != nil {
						t.Error(err)
						return
					}
				}
			}(conn)
		}
	}()

	return server
}
