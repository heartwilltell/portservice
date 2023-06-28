package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"jsonstream/app/port"
	"jsonstream/app/port/fileconsumer"
	"jsonstream/app/port/inmemstore"
	"jsonstream/pkg/log"

	"github.com/urfave/cli/v2"
)

// Variables which are related to Version command.
// Should be specified by '-ldflags' during the build phase.
// Example:
// GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Branch=$BRANCH \
// -X main.Commit=$COMMIT -o api.
var (
	// Branch is the branch this binary built from.
	Branch = "local"

	// Commit is the commit this binary built from.
	Commit = "unknown"

	// BuildTime is the time this binary built.
	BuildTime = time.Now().Format(time.RFC822)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cmd := &cli.App{
		Name:  "app",
		Usage: "root command of the app",
		Before: func(c *cli.Context) error {
			// Prints the app version.
			fmt.Printf("Branch: %s, Commit: %s, Build time: %s\n\n", Branch, Commit, BuildTime)

			return nil
		},

		Commands: []*cli.Command{
			RunCommand(),
		},
	}

	if err := cmd.RunContext(ctx, os.Args); err != nil {
		panic(err)
	}
}

func RunCommand() *cli.Command {
	cfg := struct {
		WorkerConcLimit         int
		WorkerProcessingTimeout time.Duration
		PortsFilePath           string
	}{}

	command := cli.Command{
		Name:  "run",
		Usage: "runs worker to process files with port information",
		Action: func(c *cli.Context) error {
			logger := log.New() // Init logger.

			// Initialize the storage for ports data.
			storage := inmemstore.New()

			// Initialize the data processor.
			processor := port.NewProcessor(storage)

			// Initialize the task consumer.
			consumer := fileconsumer.New(cfg.PortsFilePath)

			// Initialize the worker with all dependencies.
			worker := port.NewWorker(processor, consumer,
				port.WithLogger(logger),
				port.WithConcLimit(cfg.WorkerConcLimit),
				port.WithProcessingTimeout(cfg.WorkerProcessingTimeout),
			)

			if err := worker.Run(c.Context); err != nil && !errors.Is(err, context.Canceled) {
				return fmt.Errorf("worker failed: %w", err)
			}

			return nil
		},

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "worker-conc-limit",
				Usage:       "defines worker concurrency limit",
				EnvVars:     []string{"WORKER_CONC_LIMIT"},
				Value:       runtime.NumCPU(),
				Destination: &cfg.WorkerConcLimit,
			},
			&cli.DurationFlag{
				Name:        "worker-proc-timeout",
				Usage:       "defines worker processing timeout",
				EnvVars:     []string{"WORKER_PROC_TIMEOUT"},
				Value:       5 * time.Minute,
				Destination: &cfg.WorkerProcessingTimeout,
			},
			&cli.StringFlag{
				Name:        "ports-file-path",
				Usage:       "sets the path to json file with ports",
				Required:    true,
				Value:       "",
				Destination: &cfg.PortsFilePath,
				EnvVars:     []string{"PORTS_FILE_PATH"},
			},
		},
	}

	return &command
}
