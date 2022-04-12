package main

import (
	"context"
	"testing"
	"time"
)

func TestWorkerCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// TODO: bool -> request
	c := make(chan bool)

	go func() {
		time.Sleep(time.Second * 1)

		for i := 0; i < 10; i++ {
			c <- true
		}

		time.Sleep(time.Second * 1)
		cancel()
	}()

	startWorkers(ctx, 5, c, "./script/test.sh")
}
