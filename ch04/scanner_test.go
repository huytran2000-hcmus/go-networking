package ch04

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		message := "The bigger the interface, the weaker the abstraction"
		_, err = conn.Write([]byte(message))
		if err != nil {
			t.Error(err)
			return
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	want := []string{"The", "bigger", "the", "interface,", "the", "weaker", "the", "abstraction"}
	var got []string
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}

	if scanner.Err() != nil {
		t.Error(scanner.Err())
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
