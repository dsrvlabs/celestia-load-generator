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

func TestCallContract(t *testing.T) {
}
