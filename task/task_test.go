package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/assert"
)

const (
	nodeURL = "https://rpc.torii-1.archway.tech:443"
)

func TestRunUpload(t *testing.T) {
	ctx := context.Background()
	task := loadTask{
		ctx:     ctx,
		nodeURL: nodeURL,
		chainID: "torii-1",
		homeDir: "~/.archway",
	}

	account := "archway1jduy83242hv60p4k4kn8dfx9mv95qgrq9lrpt9"
	tx, err := task.uploadWasm(
		"../script/upload_wasm.sh",
		"../script/cw_nameservice.wasm",
		"../script/passwd",
		account,
		"torii-1",
		nodeURL,
		"~/.archway")

	fmt.Println(tx, err)
	//assert.Nil(t, err)
}

func TestGetTx(t *testing.T) {
	httpCli, err := client.NewClientFromNode(nodeURL)
	if err != nil {
		panic(err)
	}

	cliCtx := client.Context{}
	cliCtx = cliCtx.
		WithChainID("torii-1").
		WithNodeURI(nodeURL).
		WithClient(httpCli)

	tests := []struct {
		Hash         string
		ExpectHeight int64
		ExpectCode   uint32
	}{
		{
			Hash:         "50D16D73895B82666DC41FD931065CB009177A503A93F482759B6551052E81EB",
			ExpectHeight: 175880,
			ExpectCode:   0,
		},
	}

	for _, test := range tests {
		ret, err := getTx(cliCtx, test.Hash)

		assert.Nil(t, err)
		assert.Equal(t, test.ExpectCode, ret.TxResult.GetCode())
		assert.Equal(t, test.ExpectHeight, ret.Height)
	}
}

func TestUploadAndCheck(t *testing.T) {
	ctx := context.Background()
	task := loadTask{
		ctx:     ctx,
		nodeURL: nodeURL,
		chainID: "torii-1",
		homeDir: "~/.archway",
	}

	account := "archway1jduy83242hv60p4k4kn8dfx9mv95qgrq9lrpt9"
	tx, err := task.uploadWasm("../script/upload_wasm.sh", "../script/cw_nameservice.wasm", "../script/passwd", account, "torii-1", nodeURL, "~/.archway")

	assert.Nil(t, err)

	httpCli, err := client.NewClientFromNode(nodeURL)
	if err != nil {
		panic(err)
	}

	cliCtx := client.Context{}
	cliCtx = cliCtx.
		WithChainID("torii-1").
		WithNodeURI(nodeURL).
		WithClient(httpCli)

	for {
		ret, err := getTx(cliCtx, tx)
		if err != nil {
			fmt.Println("Err", err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Println(ret.TxResult.Code)
		if ret.TxResult.Code == 0 {
			fmt.Println("Success")
			break
		} else {
			fmt.Println("Failed")
			break
		}

		fmt.Println(ret.TxResult)
		time.Sleep(1 * time.Second)
	}
}

/*
func TestTask(t *testing.T) {

	sChan := make(chan int, 1)
	fChan := make(chan int, 1)

	go func() {
		for {
			select {
			case <-sChan:
				log.Println("S")
			case <-fChan:
				log.Println("F")
			}
		}
	}()

	accounts := []string{
		"archway1jduy83242hv60p4k4kn8dfx9mv95qgrq9lrpt9",
		"archway1pm0yyd2ncc2x67ctuz5p3tcxa59tezx5scp0hj",
		"archway179jgnt5ckmnjrxg2rykycv4vkh2y8nh8h30sga",
		"archway1fqdch0dl4r43wp3yw6f5nkp7jca3xn4m3s2plh",
		"archway13yazsem53w0lpsf0j0n0l7038jvs5xadr85zrt",
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	_ = cancel

	task := NewLoadTask(ctx)
	task.StartUpload(accounts, sChan, fChan)
}

*/
