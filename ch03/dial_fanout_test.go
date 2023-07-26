package ch03

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	dial := func(ctx context.Context, id int, res chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()

		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", listener.Addr().String())
		if err != nil {
			return
		}
		conn.Close()

		select {
		case <-ctx.Done():
		case res <- id:
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	res := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, i, res, &wg)
	}

	id := <-res
	cancel()
	wg.Wait()
	close(res)

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Errorf("want cancelled error, got: %v", ctx.Err())
	}

	t.Logf("dialer %d retrieved the result", id)
}
