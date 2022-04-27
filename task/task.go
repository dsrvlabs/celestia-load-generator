package task

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"os/exec"
	"strings"
	"sync"
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
	StartCall(accounts []string, address string, sChan, fChan chan<- int) error
}

type loadTask struct {
	ctx     context.Context
	chainID string
	nodeURL string
	homeDir string
}

func (t *loadTask) StartUpload(accounts []string, wasmFile, passwdFile string, sChan, fChan chan<- int) error {
	wg := sync.WaitGroup{}

	for _, acc := range accounts {
		go t.taskUpload(acc, wasmFile, passwdFile, sChan, fChan, &wg)
	}

	for {
		select {
		case <-t.ctx.Done():
			return nil
		}
	}
}

func (t *loadTask) taskUpload(account, wasmFile, passwdFile string, sChan, fChan chan<- int, wg *sync.WaitGroup) {
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
		txHash, err := t.uploadWasm("./upload_wasm.sh", wasmFile, passwdFile, account, t.chainID, t.nodeURL, t.homeDir)
		if err != nil {
			log.Println("uploadWasm", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("Hash", txHash)

		for {
			ret, err := getTx(cliCtx, txHash)
			if err != nil {
				log.Println("Err", err)
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

func (t *loadTask) uploadWasm(filename string, args ...string) (string, error) {
	log.Println("uploadWasm", args)

	cmd := exec.Command(filename, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Println("Err uploadWasm", err)
		log.Println("uploadWasm Out", out.String())
		return "", err
	}

	log.Println("uploadWasm Out", out.String())

	txHash := strings.Trim(out.String(), "\r\n")

	log.Println("uploadWasm Hash", txHash)

	return txHash, nil
}

func getTx(ctx client.Context, hash string) (*ctypes.ResultTx, error) {
	log.Println("Call getTx", hash)

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

func (t *loadTask) StartCall(accounts []string, address string, sChan, fChan chan<- int) error {
	// TODO: TBD
	return nil
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
