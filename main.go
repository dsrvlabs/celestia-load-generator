package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rootwarp/wasm-load-tester/task"
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
			fmt.Println("Args", args)

			flags := cmd.Flags()

			wasmFile, err := flags.GetString("wasm")
			if err != nil {
				log.Println(err)
				return err
			}

			passwdFile, err := flags.GetString("password")
			if err != nil {
				log.Println(err)
				return err
			}

			accountFile, err := flags.GetString("account")
			if err != nil {
				log.Println(err)
				return err
			}

			chainID, err := flags.GetString("chain-id")
			if err != nil {
				log.Println(err)
				return err
			}

			nodeURL, err := flags.GetString("node")
			if err != nil {
				log.Println(err)
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

				_ = <-c

				log.Println("Interrupt")
				cancel()
			}()

			// TODO:
			loader := task.NewLoadTask(ctx, chainID, nodeURL, "~/.archway")

			f, err := os.Open(accountFile)
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

			loader.StartUpload(accounts, wasmFile, passwdFile, sChan, fChan)

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

	callCmd := &cobra.Command{
		Use: "call",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(uploadCmd)
	cmd.AddCommand(callCmd)

	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}

func startLoader(ctx context.Context, workers int, loadCmd string) {
	sChan := make(chan int, channelBuffer)
	fChan := make(chan int, channelBuffer)
	cmdChan := make(chan string, channelBuffer)

	statChan := tpsCalculator(ctx, sChan, fChan)

	go addQueue(ctx, cmdChan, loadCmd)
	go printTPS(ctx, statChan)
	startWorkers(ctx, workers, cmdChan, sChan, fChan)
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
