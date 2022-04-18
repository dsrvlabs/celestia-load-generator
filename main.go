package main

import (
	"context"
	"log"
	"os/exec"
	"sync"
	"time"
)

const (
	elapseTickPeriod = 100 * time.Millisecond
)

type statistic struct {
	Time           time.Time
	TotalRequest   int64
	SuccessRequest int64
	FailRequest    int64
	Elapse         time.Duration
}

// TODO: Get N parameter.
// TODO: Target TPS.
func main() {
	sChan := make(chan int, 100)
	fChan := make(chan int, 100)
	cmdChan := make(chan string, 100)

	statChan := tpsCalculator(sChan, fChan)

	go addQueue(cmdChan)
	go printTPS(statChan)

	ctx := context.Background()
	startWorkers(ctx, 10, cmdChan, sChan, fChan)
}

func addQueue(cmdChan chan<- string) {
	for {
		_ = <-time.Tick(1 * time.Millisecond)
		cmdChan <- "./script/test.sh"
	}

}

func printTPS(statChan <-chan statistic) {
	for {
		stat := <-statChan
		log.Printf("Req %d/%d, TPS: %f\n", stat.SuccessRequest, stat.TotalRequest, float64(stat.SuccessRequest)/stat.Elapse.Seconds())
	}
}

func startWorkers(ctx context.Context, n int, cmdChan <-chan string, successChan, failChan chan<- int) {
	log.Println("Start Worker")

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(idx int) {
			log.Println("Start worker", idx)
			defer log.Println("Stop worker")
			defer wg.Done()

			for {
				select {
				case cmd, ok := <-cmdChan:
					if !ok {
						return
					}

					c := exec.Command(cmd)
					err := c.Run()
					if err != nil {
						failChan <- 1
					} else {
						successChan <- 1
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	wg.Wait()
}

func tpsCalculator(successChan, failChan <-chan int) chan statistic {
	tpsChan := make(chan statistic, 1)

	go func() {
		stat := statistic{}

		for {
			select {
			case n := <-successChan:
				stat.SuccessRequest += int64(n)
				stat.TotalRequest += int64(n)
			case n := <-failChan:
				stat.FailRequest += int64(n)
				stat.TotalRequest += int64(n)
			case <-time.Tick(elapseTickPeriod):
				stat.Time = time.Now()
				stat.Elapse += elapseTickPeriod
				tpsChan <- stat
			}
		}
	}()

	return tpsChan
}

// TODO: Load wasm binary by code.
// TODO: Call wasm contract. Require heavy load contract.
// TODO: Report output.
