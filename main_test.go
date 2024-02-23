package main

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

// testServer runs a server to test against
func testServer(ctx context.Context) {
	ln, err := new(net.ListenConfig).Listen(ctx, "tcp", Address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		switch {
		case errors.Is(err, net.ErrClosed):
			return
		case err != nil:
			panic(err)
		}

		// handles unlimited connections
		go func() {
			time.Sleep(time.Millisecond * 100)
			io.Copy(conn, strings.NewReader("hello world"))
			conn.Close()
		}()
	}
}
