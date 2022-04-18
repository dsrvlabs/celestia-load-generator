package main

import (
	"context"
	"testing"
	"time"
)

func TestWorkerCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan string)

	go func() {
		time.Sleep(time.Second * 1)

		for i := 0; i < 10; i++ {
			c <- "./script/test.sh"
		}

		cancel()
	}()

	successChan := make(chan int)
	failChan := make(chan int)

	startWorkers(ctx, 5, c, successChan, failChan)
}

func TestTPSCalculation(t *testing.T) {
	sChan := make(chan int, 100)
	fChan := make(chan int, 100)
	cmdChan := make(chan string, 100)

	statChan := tpsCalculator(sChan, fChan)

	go addQueue(cmdChan)
	go printTPS(statChan)

	ctx := context.Background()
	startWorkers(ctx, 10, cmdChan, sChan, fChan)
}
