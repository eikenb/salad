package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestFetchMsg(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener := testServer(testMessages...)

	addr := listener.Addr().String()
	for _, want := range testMessages {
		got, err := fetchMsg(ctx, addr)
		if err != nil {
			t.Fatalf("unexpected error: `%s`", err)
		}
		if want != *got {
			t.Errorf("message parse error; want: `%#v`, got: `%#v`", want, got)
		}
	}
}

func TestBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origBackoff := baseBackoff
	baseBackoff = time.Nanosecond
	defer func() { baseBackoff = origBackoff }()

	t.Run("backoff increment", func(t *testing.T) {
		attempt := 0
		for want := range maxBackoff {
			if attempt != want {
				t.Errorf("backoff not incrementing properly")
			}
			attempt = backoff(ctx, attempt) // should increment by 1
		}
	})

	t.Run("max backoff", func(t *testing.T) {
		attempt := 0
		want := maxBackoff
		for range maxBackoff + 10 {
			attempt = backoff(ctx, attempt) // should increment by 1
		}
		if attempt != want {
			t.Errorf("not maxing out; got: %d, want: %d", attempt, want)
		}
	})
}

// testServer runs a server to test against
func testServer(tests ...message) net.Listener {
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
	}()
	return ln
}
