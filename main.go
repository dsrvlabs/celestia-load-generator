package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dsrvlabs/wasm-load-generator/task"
	"github.com/spf13/cobra"
)

const (
	elapseTickPeriod   = 100 * time.Millisecond
	channelBuffer      = 100
	defaultWorkerCount = 10
)

type statistic struct {
	Time           time.Time
	TotalRequest   int64
	SuccessRequest int64
	FailRequest    int64
	Elapse         time.Duration
}

type flagConfigs struct {
	WasmFile       string `json:"wasm_file"`
	PasswordFile   string `json:"password_file"`
	AccountFile    string `json:"account_file"`
	ChainID        string `json:"chain_id"`
	Node           string `json:"node"`
	ContractAddess string `json:"contract_addess"`
}

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("archway", "archwaypub")
	cfg.SetBech32PrefixForValidator("archwayvaloper", "archwayvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("archwayvalcons", "archwayvalconspub")
	cfg.Seal()
}

func main() {
	cmd := cobra.Command{}

	uploadCmd := &cobra.Command{
		Use: "upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags, err := parseFlags(cmd)

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

				_ = <-c

				log.Println("Interrupt")
				cancel()
			}()

			// TODO: Home dir
			loader := task.NewLoadTask(ctx, flags.ChainID, flags.Node, "~/.archway")

			f, err := os.Open(flags.AccountFile)
			if err != nil {
				log.Panic(err)
			}

			r := bufio.NewReader(f)

			accounts := []string{}
			for {
				line, _, err := r.ReadLine()
				if err != nil {
					break
				}

				accounts = append(accounts, string(line))
			}

			sChan := make(chan int, channelBuffer)
			fChan := make(chan int, channelBuffer)
			statChan := tpsCalculator(ctx, sChan, fChan)

			go printTPS(ctx, statChan)

			loader.StartUpload(accounts, flags.WasmFile, flags.PasswordFile, sChan, fChan)

			return nil
		},
	}

	uploadCmd.Flags().StringP("wasm", "w", "", "WASM file")
	uploadCmd.MarkFlagRequired("wasm")

	uploadCmd.Flags().StringP("password", "p", "", "Password file")
	uploadCmd.MarkFlagRequired("password")

	uploadCmd.Flags().StringP("account", "a", "", "account file")
	uploadCmd.MarkFlagRequired("account")

	uploadCmd.Flags().StringP("chain-id", "c", "", "chain id")
	uploadCmd.MarkFlagRequired("chain-id")

	uploadCmd.Flags().StringP("node", "n", "", "Node ID")
	uploadCmd.MarkFlagRequired("node")

	// TODO: Refactoring this.
	callCmd := &cobra.Command{
		Use: "call",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags, err := parseFlags(cmd)

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

				_ = <-c

				log.Println("Interrupt")
				cancel()
			}()

			// TODO: Home dir
			loader := task.NewLoadTask(ctx, flags.ChainID, flags.Node, "~/.archway")

			f, err := os.Open(flags.AccountFile)
			if err != nil {
				log.Panic(err)
			}

			r := bufio.NewReader(f)

			accounts := []string{}
			for {
				line, _, err := r.ReadLine()
				if err != nil {
					break
				}

				accounts = append(accounts, string(line))
			}

			sChan := make(chan int, channelBuffer)
			fChan := make(chan int, channelBuffer)
			statChan := tpsCalculator(ctx, sChan, fChan)

			go printTPS(ctx, statChan)

			loader.StartCall(accounts, flags.PasswordFile, flags.ContractAddess, sChan, fChan)

			return nil
		},
	}

	callCmd.Flags().StringP("password", "p", "", "Password file")
	callCmd.MarkFlagRequired("password")

	callCmd.Flags().StringP("account", "a", "", "account file")
	callCmd.MarkFlagRequired("account")

	callCmd.Flags().StringP("contract", "t", "", "contract address")
	callCmd.MarkFlagRequired("contract")

	callCmd.Flags().StringP("chain-id", "c", "", "chain id")
	callCmd.MarkFlagRequired("chain-id")

	callCmd.Flags().StringP("node", "n", "", "Node ID")
	callCmd.MarkFlagRequired("node")

	cmd.AddCommand(uploadCmd)
	cmd.AddCommand(callCmd)

	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}

func parseFlags(cmd *cobra.Command) (*flagConfigs, error) {
	retFlags := &flagConfigs{}

	flags := cmd.Flags()
	if wasmFile, err := flags.GetString("wasm"); err == nil {
		wasmFile, err = filepath.Abs(wasmFile)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		retFlags.WasmFile = wasmFile
	}

	if passwdFile, err := flags.GetString("password"); err == nil {
		passwdFile, err = filepath.Abs(passwdFile)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		retFlags.PasswordFile = passwdFile
	}

	if accountFile, err := flags.GetString("account"); err == nil {
		accountFile, err = filepath.Abs(accountFile)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		retFlags.AccountFile = accountFile
	}

	if contract, err := flags.GetString("contract"); err == nil {
		log.Println(err)
		retFlags.ContractAddess = contract
		return nil, err
	}

	chainID, err := flags.GetString("chain-id")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	retFlags.ChainID = chainID

	nodeURL, err := flags.GetString("node")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	retFlags.Node = nodeURL

	return retFlags, nil
}

func tpsCalculator(ctx context.Context, successChan, failChan <-chan int) chan statistic {
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
			case <-ctx.Done():
				break
			}
		}
	}()

	return tpsChan
}

func addQueue(ctx context.Context, cmdChan chan<- string, loadCmd string) {
	for {
		select {
		case <-time.Tick(1 * time.Millisecond):
			cmdChan <- loadCmd
		case <-ctx.Done():
			break
		}
	}
}

func printTPS(ctx context.Context, statChan <-chan statistic) {
	for {
		select {
		case stat := <-statChan:
			log.Printf("Req %d/%d, TPS: %f\n", stat.SuccessRequest, stat.TotalRequest, float64(stat.SuccessRequest)/stat.Elapse.Seconds())
			continue
		case <-ctx.Done():
			break
		}
	}
}

// TODO: Load wasm binary by code.
// TODO: Call wasm contract. Require heavy load contract.
// TODO: Report output.
