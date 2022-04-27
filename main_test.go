package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
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

	startLoader(ctx, 10, "")
}

// func TestUploadCommand(t *testing.T) {
// 	cmd := exec.Command(
// 		"/home/rootwarp/bin/archwayd",
// 		"tx", "wasm", "store",
// 		"./script/cw_nameservice.wasm",
// 		"--chain-id=torii-1",
// 		"--from=cheese",
// 		"--home=~/.archway",
// 		"--gas=auto",
// 		"--broadcast-mode=sync",
// 		"--node=https://rpc.torii-1.archway.ech:443",
// 		"-y",
// 		"--output=json",
// 	)
//
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &out
// 	cmd.Stdin = strings.NewReader("@validator2022")
//
// 	err := cmd.Run()
//
// 	fmt.Println(err)
// 	fmt.Println(out.String())
// }

func TestExec(t *testing.T) {
	cmd := exec.Command(
		"bash",
		"./upload_wasm.sh",
		"./cw_nameservice.wasm",
		"./passwd",
		"cheese",
		"torii-1",
		"https://rpc.torii-1.archway.tech:443",
		"~/.archway",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	fmt.Println(err)
	fmt.Println(out.String())
}
