package fileconsumer

import (
	"context"
)

// Consumer implements port.Consumer interface by doing nothing.
type Consumer struct{ files []string }

// New returns a pointer to a new instance of Consumer.
func New(files ...string) *Consumer {
	c := Consumer{
		files: append(make([]string, 0, len(files)), files...),
	}

	return &c
}

func (c *Consumer) Close() error { return nil }

func (c *Consumer) Consume(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	var file string

	if len(c.files) > 0 {
		file = c.files[0]
		c.files = c.files[1:]
	}

	return file, nil
}
