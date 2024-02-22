package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// hardcoded address to keep things simple
// const Address = "data.salad.com:5000"
const Address = "localhost:5000"

// main loop
func main() {
	ctx := SignalContext()

	go testServer(ctx)

	d := new(net.Dialer)
	for {
		// connect to server via tcp
		conn, err := d.DialContext(ctx, "tcp", Address)
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			return // ie. exit (dialer's context cancelled via signal)
		case err != nil:
			panic(err)
		}
		defer conn.Close()
		// wait for data
		//   retry on conn errors
		// receive data in buffer
		//   retry on read errors
		buf, _ := io.ReadAll(conn)
		// retry if error before all read

		// parse data into struct
		// pretty print struct
		fmt.Println(string(buf))
		// loop
	}
}

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
