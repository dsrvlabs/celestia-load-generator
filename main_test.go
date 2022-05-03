package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestAbsolutePath(t *testing.T) {
	absPath, err := filepath.Abs("./README.md")

	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(absPath, "/"))
	assert.True(t, strings.HasSuffix(absPath, "/README.md"))
}
