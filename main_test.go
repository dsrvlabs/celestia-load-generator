package main

import (
	"context"
	"testing"
	"time"
)

func TestTPSCalculation(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	startLoader(ctx, 10)
}
