package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"testing"
)

func TestDialLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener := testServer(cancel, testMessages...)

	results := make(chan message)
	output := func(msg message) {
		results <- msg
	}

	addr := listener.Addr().String()
	go dialLoop(ctx, addr, output)
	for _, want := range testMessages {
		got := <-results
		if want != got {
			t.Errorf("message parse error; want: `%#v`, got: `%#v`", want, got)
		}
	}
}

// testServer runs a server to test against
func testServer(cancel func(), tests ...message) net.Listener {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	go func() {
		defer ln.Close()
		for _, msg := range tests {
			conn, err := ln.Accept()
			switch {
			case errors.Is(err, context.Canceled):
				return
			case err != nil:
				panic(err)
			}
			bin, err := msg.MarshalBinary()
			if err != nil {
				panic(err)
			}
			io.Copy(conn, bytes.NewReader(bin))
			conn.Close()
		}
		cancel()
	}()
	return ln
}
