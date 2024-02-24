package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

// hardcoded address to keep things simple
// const Address = "data.salad.com:5000"
var Address = "localhost:5000"

// inject server address into process environment for production config
func init() {
	address := os.Getenv("SERVER_ADDRESS")
	if address != "" {
		Address = address
	}
}

func main() {
	ctx := SignalContext()
	output := func(msg message) {
		fmt.Printf("%+v\n", msg)
	}

	dialLoop(ctx, Address, output)
}

// dialLoop runs the main dial loop connecting to and reading in the data from
// the server. It is split out for testing.
func dialLoop(ctx context.Context, addr string, output func(message)) {
	d := new(net.Dialer)
	buf := new(bufio.Reader)
	attempt := 0
	for {
		// connect to server via tcp
		conn, err := d.DialContext(ctx, "tcp", addr)
		//   retry on conn errors
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			return // ie. exit (dialer's context cancelled via signal)
		case err != nil:
			slog.Info("DialContext", "error", err)
			attempt = backoff(attempt) // very simple exponential backoff
			continue
		}
		// wait for data
		// receive data in buffer
		buf.Reset(conn)
		data, err := io.ReadAll(buf)
		// retry if error before all read
		if err != nil {
			slog.Info("ReadAll", "error", err)
			continue
		}
		// parse data into struct
		msg := new(message)
		msg.UnmarshalBinary(data)
		// pretty print struct
		output(*msg)
		// end loop, reset attempts and close this connection
		attempt = 0
		conn.Close()
	}
}

// retry with backoff
func backoff(attempt int) int {
	wait := time.Second * time.Duration(attempt*attempt)
	time.Sleep(wait)
	return attempt + 1
}
