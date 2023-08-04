package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	count    = flag.Int("c", 3, "number of pings")
	interval = flag.Duration("i", 1*time.Second, "time in second between pings")
	timeout  = flag.Duration("t", 3*time.Second, "time to wait for reply")
)

func init() {
	flag.Usage = func() {
		fmt.Printf(`Usage: %s [options] host:port
Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("host:port is required\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *count <= 0 {
		fmt.Printf("invalid number of pings: %d\n", *count)
		os.Exit(1)
	}

	target := flag.Arg(0)
	fmt.Println("PING:", target)

	for nMsg := 0; nMsg < *count; nMsg++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, *timeout)
		duration := time.Since(start)
		if err != nil {
			var nErr net.Error
			ok := errors.As(err, &nErr)
			if !ok && !nErr.Timeout() {
				fmt.Printf("encouter error: %s\n", err)
				os.Exit(1)
			}
		} else {
			_ = conn.Close()
			fmt.Printf("Duration: %v\n", duration)
		}

		time.Sleep(*interval)
	}
}
