package main

import (
	"context"
	"testing"
	"time"
)

func TestTPSCalculation(t *testing.T) {
	sChan := make(chan int, 100)
	fChan := make(chan int, 100)
	cmdChan := make(chan string, 100)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	statChan := tpsCalculator(ctx, sChan, fChan)

	go addQueue(ctx, cmdChan)
	go printTPS(ctx, statChan)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	startWorkers(ctx, 10, cmdChan, sChan, fChan)
}
