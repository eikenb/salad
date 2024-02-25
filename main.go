package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"os"
	"time"
)

var Address = "data.salad.com:5000"
var maxBackoff = 7 // ~5 sec max w/ 100ms base
var baseBackoff = time.Millisecond * 100

var retryErr = errors.New("retry")

// inject server address into process environment for production config
func init() {
	address := os.Getenv("SERVER_ADDRESS")
	if address != "" {
		Address = address
	}
}

// main loop containing as little of the logic as possible.
func main() {
	ctx := SignalContext()
	attempt := 0
	for {
		msg, err := fetchMsg(ctx, Address)
		switch err {
		case context.Canceled:
			return
		case retryErr:
			attempt = backoff(ctx, attempt)
			continue
		}
		attempt = 0
		fmt.Printf("%+v\n", msg)
	}
}

// backoff implements a simple exponential backoff
func backoff(ctx context.Context, attempt int) int {
	if attempt >= maxBackoff {
		attempt = maxBackoff - 1
	}
	jitter := time.Duration(rand.Float64() * float64(baseBackoff))
	wait := (baseBackoff * time.Duration(attempt*attempt)) + jitter
	select {
	case <-time.After(wait):
	case <-ctx.Done():
	}
	return attempt + 1
}

// fetchMsg dials the server, grabs one binary blob, converts it to a message
// and returns it or an error. If an error it will use the `retryErr` and
// context.Canceled errors to communicate to the loop.
func fetchMsg(ctx context.Context, addr string) (*message, error) {
	var d net.Dialer
	// connect to server via tcp
	conn, err := d.DialContext(ctx, "tcp", addr)
	//   retry on conn errors
	switch {
	case errors.Is(ctx.Err(), context.Canceled):
		return nil, ctx.Err() // ie. exit (dialer's context cancelled via signal)
	case err != nil:
		slog.Info("DialContext", "error", err)
		return nil, retryErr
	}
	defer conn.Close()
	// wait for data
	data, err := io.ReadAll(conn)
	// retry if error before all read
	if err != nil {
		slog.Info("ReadAll", "error", err)
		return nil, retryErr
	}
	// parse data into struct
	msg := new(message)
	msg.UnmarshalBinary(data)
	return msg, nil
}
