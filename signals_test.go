package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_signalHandling(t *testing.T) {
	t.Run("signal should exit", func(t *testing.T) {
		signalFaker := make(chan os.Signal, 1)
		signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
		for _, sig := range signals {
			signalFaker <- sig
			select {
			case <-signalHandling(signalFaker).Done():
			case <-time.After(time.Second):
				t.Error("should not have blocked")
			}
		}
	})
	t.Run("channel close should exit", func(t *testing.T) {
		signalFaker := make(chan os.Signal)
		ctx := signalHandling(signalFaker)
		select {
		case <-ctx.Done():
			t.Error("should not return")
		default:
		}
		select {
		case signalFaker <- syscall.SIGINT:
		default:
			t.Error("signal should have sent")
		}
		select {
		case signalFaker <- syscall.SIGINT:
			t.Error("signal should NOT have sent")
		default:
		}
	})
}
