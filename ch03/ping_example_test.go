package ch03

import (
	"context"
	"fmt"
	"io"
	"time"
)

func ExamplePINGer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, w := io.Pipe()

	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second

	done := make(chan struct{})
	go func() {
		go Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("resetting timer (%s)\n", d)
			resetTimer <- d
		}

		buf := make([]byte, 1024)
		start := time.Now()
		n, err := r.Read(buf)
		end := time.Since(start)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("received %q (%s)\n", buf[:n], end.Round(100*time.Millisecond))
	}

	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		fmt.Printf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	<-done
	// Output:
	// Run 1:
	// resetting timer (0s)
	// received "PING" (1s)
	// Run 2:
	// resetting timer (200ms)
	// received "PING" (200ms)
	// Run 3:
	// resetting timer (300ms)
	// received "PING" (300ms)
	// Run 4:
	// resetting timer (0s)
	// received "PING" (300ms)
	// Run 5:
	// received "PING" (300ms)
	// Run 6:
	// received "PING" (300ms)
	// Run 7:
	// received "PING" (300ms)
}
