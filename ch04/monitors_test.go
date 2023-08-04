package ch04

import (
	"io"
	"log"
	"net"
	"os"
)

func ExampleMonitor() {
	monitor := Monitor{Logger: log.New(os.Stdout, "monitor: ", log.LstdFlags)}
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		conn, err := listener.Accept()
		if err != nil {
			monitor.Fatal(err)
			return
		}
		defer conn.Close()

		r := io.TeeReader(conn, monitor)

		b := make([]byte, 256)
		n, err := r.Read(b)
		if err != nil && err != io.EOF {
			monitor.Fatal(err)
			return
		}

		w := io.MultiWriter(conn, monitor)

		_, err = w.Write(b[:n])
		if err != nil {
			monitor.Fatal(err)
			return
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte("Testing..."))
	if err != nil {
		monitor.Fatal(err)
	}
}
