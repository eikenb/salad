package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type empty = struct{}

// SignalContext returns a context that is cancelled when the process receives
// a SIGINT or SIGTERM
func SignalContext() context.Context {
	signals := make(chan os.Signal, 1)
	return signalHandling(signals)
}

// signalHandling takes a signals channel and returns a context that is
// cancelled upon receiving a SIGINT or SIGTERM.
// This is split out for testing. Use SignalContext.
func signalHandling(signals chan os.Signal) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	wait := make(chan empty)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func(signals chan os.Signal) {
		close(wait)
		for {
			select {
			case sig := <-signals:
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM:
					cancel()
					signal.Reset()
					return
				}
			}
		}
	}(signals)
	<-wait // wait for goroutine to start before returning, eases testing
	return ctx
}
