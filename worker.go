package main

import (
	"context"
	"log"
	"os/exec"
	"sync"
)

func startWorkers(ctx context.Context, n int, cmdChan <-chan string, successChan, failChan chan<- int) {
	log.Println("Start Worker")

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(idx int) {
			log.Println("Start worker", idx)
			defer log.Println("Stop worker")
			defer wg.Done()

			key, err := createKey()
			if err != nil {
				panic(err)
			}

			log.Println("New key created", key.Address)

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
