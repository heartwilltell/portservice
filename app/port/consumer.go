package port

import (
	"context"
	"io"
)

// Consumer wraps the logic of receiving the task that should be processed.
type Consumer interface {
	io.Closer

	// Consume consumes the task from the remote source and returns a path
	// to the file that should be processed by the Processor.
	Consume(ctx context.Context) (string, error)
}
