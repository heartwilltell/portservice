package port

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"jsonstream/pkg/log"
)

// WorkerOption configures the Worker.
type WorkerOption func(worker *Worker)

// WithLogger sets the given logger as Worker logger.
func WithLogger(logger log.Logger) WorkerOption {
	return func(w *Worker) { w.logger = logger }
}

// WithConcLimit sets the given limit as Worker concLimit.
func WithConcLimit(limit int) WorkerOption {
	return func(w *Worker) { w.concLimit = limit }
}

func WithProcessingTimeout(timeout time.Duration) WorkerOption {
	return func(w *Worker) { w.procTimeout = timeout }
}

// Worker represents a worker process that consumes a tasks
// via Consumer and passes it to the Processor for processing.
type Worker struct {
	logger log.Logger

	concLimit   int
	procTimeout time.Duration

	consumer  Consumer
	processor Processor
}

// NewWorker returns a pointer to a new instance of Worker.
func NewWorker(processor Processor, consumer Consumer, options ...WorkerOption) *Worker {
	w := Worker{
		logger: log.DisabledLogger(),

		// Num of logical cores.
		concLimit:   runtime.NumCPU(),
		procTimeout: 5 * time.Minute,

		consumer:  consumer,
		processor: processor,
	}

	for _, option := range options {
		option(&w)
	}

	return &w
}

func (w *Worker) Run(ctx context.Context) error {
	tasks := make(chan string, w.concLimit)

	var wg sync.WaitGroup
	wg.Add(w.concLimit)

	for i := 0; i < w.concLimit; i++ {
		w.process(&wg, tasks)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Infof("Shutdown signal received: %s", ctx.Err().Error())

			// Stop consuming the tasks.
			if err := w.consumer.Close(); err != nil {
				return fmt.Errorf("closing consumer: %w", err)
			}

			// Close task channel to stop all processor routines.
			close(tasks)

			// Wait until all worker will finish the ongoing processing.
			wg.Wait()

			return ctx.Err()

		default:
			task, err := w.consumer.Consume(ctx)
			if err != nil {
				return fmt.Errorf("consuming task: %w", err)
			}

			if task != "" {
				tasks <- task
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (w *Worker) process(wg *sync.WaitGroup, tasks <-chan string) {
	defer wg.Done()

	for task := range tasks {
		ctx, cancel := context.WithTimeout(context.Background(), w.procTimeout)

		if err := w.processor.Process(ctx, task); err != nil {
			// TODO: consider to send an error to error tracking system (Sentry.io).
			w.logger.Errorf("Failed to process '%s' task: %s", task, err.Error())
		} else {
			w.logger.Infof("File %s has been processed", task)
		}

		cancel()
	}
}
