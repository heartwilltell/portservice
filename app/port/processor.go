package port

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Port represents an information about a concrete port.
type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
}

// Processor wraps the logic of processing ports data.
type Processor interface {
	// Process tries to read and process batch file with an information about ports.
	Process(ctx context.Context, filePath string) error
}

// Storage represents a persistence layer logic.
type Storage interface {
	// WritePort writes an information about port to the database.
	WritePort(ctx context.Context, id string, port *Port) error
}

// JSONFileProcessor implements Processor interface
// for json files which it reads as a single batch.
type JSONFileProcessor struct {
	storage Storage
}

// NewProcessor returns a pointer to a new instance of JSONFileProcessor.
func NewProcessor(storage Storage) *JSONFileProcessor {
	return &JSONFileProcessor{storage: storage}
}

func (p *JSONFileProcessor) Process(ctx context.Context, filePath string) (bErr error) {
	f, openErr := os.Open(filePath)
	if openErr != nil {
		return fmt.Errorf("json file streamer: open '%s' file: %w", filePath, openErr)
	}

	defer func() {
		if err := f.Close(); err != nil {
			bErr = errors.Join(bErr, fmt.Errorf("closing '%s' file: %w", filePath, err))
		}
	}()

	reader := bufio.NewReader(f)
	decoder := json.NewDecoder(reader)

	if err := p.handleOpenCloseDelimiter(decoder); err != nil {
		return fmt.Errorf("parsing open delimiter: %w", err)
	}

	for decoder.More() {
		// Check for context cancellation on each iteration.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			return fmt.Errorf("decode port id token: %w", tokenErr)
		}

		id, ok := token.(string)
		if !ok {
			return fmt.Errorf("decode port id: not a string")
		}

		// TODO: consider to add sync.Pool to avoid additional allocation on each iteration.
		var port *Port

		if err := decoder.Decode(&port); err != nil {
			return fmt.Errorf("decode port: %w", err)
		}

		if err := p.storage.WritePort(ctx, id, port); err != nil {
			return fmt.Errorf("save port to storage: %w", err)
		}
	}

	if err := p.handleOpenCloseDelimiter(decoder); err != nil {
		return fmt.Errorf("parsing close delimiter:: %w", err)
	}

	return nil
}

func (*JSONFileProcessor) handleOpenCloseDelimiter(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("handle open/close delimiter: %w", err)
	}

	delimiter, ok := token.(json.Delim)
	if !ok {
		return fmt.Errorf("invalid json delimiter: %v", delimiter)
	}

	if !strings.ContainsAny(string(delimiter), "{[]}") {
		return fmt.Errorf("invalid json delimiter: %s", delimiter)
	}

	return nil
}
