package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

// TODO: Get N parameter.
// TODO: Target TPS.
func main() {
	fmt.Println("vim-go")
}

func startWorkers(ctx context.Context, n int, c <-chan bool, cmd string) {
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
				case r, ok := <-c:
					if !ok {
						return
					}

					log.Println("Handle", idx)
					cmd := exec.Command(cmd)
					err := cmd.Run()

					log.Println(err)
					_ = r
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	wg.Wait()
}

// TODO: Load wasm binary by code.
// TODO: Call wasm contract. Require heavy load contract.
// TODO: Show TPS.
