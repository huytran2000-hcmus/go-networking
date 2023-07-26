package ch03

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	interval := defaultPingInterval
	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			_, err := w.Write([]byte("PING"))
			if err != nil {
				return
			}
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		}
		_ = timer.Reset(interval)
	}
}
