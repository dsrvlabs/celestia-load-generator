package task

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type uploadResponse struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
}

// LoadTask provides interfaces for load testing.
type LoadTask interface {
	StartUpload(accounts []string, wasmFile, passwdFile string, sChan, fChan chan<- int) error

	StartCall(accounts []string, passwdFile, address string, sChan, fChan chan<- int) error
}

type loadTask struct {
	ctx     context.Context
	chainID string
	nodeURL string
	homeDir string
}

func (t *loadTask) StartUpload(accounts []string, wasmFile, passwdFile string, sChan, fChan chan<- int) error {
	for _, acc := range accounts {
		go t.taskUpload(acc, wasmFile, passwdFile, sChan, fChan)
	}

	for {
		select {
		case <-t.ctx.Done():
			return nil
		}
	}
}

func (t *loadTask) taskUpload(account, wasmFile, passwdFile string, sChan, fChan chan<- int) {
	httpCli, err := client.NewClientFromNode(t.nodeURL)
	if err != nil {
		panic(err)
	}

	cliCtx := client.Context{}
	cliCtx = cliCtx.
		WithChainID(t.chainID).
		WithNodeURI(t.nodeURL).
		WithClient(httpCli)

	for {
		// TODO: Script path.
		txHash, err := t.executeShellScript("./upload_wasm.sh", wasmFile, passwdFile, account, t.chainID, t.nodeURL, t.homeDir)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			ret, err := getTx(cliCtx, txHash)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if ret.TxResult.Code == 0 {
				log.Println("Success")
				sChan <- 1
				break
			} else {
				log.Println("Failed")
				fChan <- 1
				break
			}
		}
	}
}

func (t *loadTask) executeShellScript(filename string, args ...string) (string, error) {
	cmd := exec.Command(filename, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	txHash := strings.Trim(out.String(), "\r\n")

	return txHash, nil
}

func getTx(ctx client.Context, hash string) (*ctypes.ResultTx, error) {
	h, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
		//return nil, err
	}

	cli, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	return cli.Tx(context.Background(), h, true)
}

func (t *loadTask) StartCall(accounts []string, passwdFile, address string, sChan, fChan chan<- int) error {
	for _, acc := range accounts {
		go t.taskContractCall(acc, passwdFile, address, sChan, fChan)
	}

	for {
		select {
		case <-t.ctx.Done():
			return nil
		}
	}
}

func (t *loadTask) taskContractCall(account, passwdFile, address string, sChan, fChan chan<- int) {
	httpCli, err := client.NewClientFromNode(t.nodeURL)
	if err != nil {
		panic(err)
	}

	cliCtx := client.Context{}
	cliCtx = cliCtx.
		WithChainID(t.chainID).
		WithNodeURI(t.nodeURL).
		WithClient(httpCli)

	for {
		//
		txHash, err := t.executeShellScript("./call_contract.sh", passwdFile, account, address, t.chainID, t.nodeURL, t.homeDir)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			ret, err := getTx(cliCtx, txHash)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if ret.TxResult.Code == 0 {
				log.Println("Success")
				sChan <- 1
				break
			} else {
				log.Println("Failed")
				fChan <- 1
				break
			}
		}
	}
}

// NewLoadTask creates new test instances.
func NewLoadTask(ctx context.Context, chainID, nodeURL, homeDir string) LoadTask {
	return &loadTask{
		ctx:     ctx,
		nodeURL: nodeURL,
		chainID: chainID,
		homeDir: homeDir,
	}
}
