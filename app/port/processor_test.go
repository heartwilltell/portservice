package port

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestJsonFileProcessor_Process(t *testing.T) {
	type tcase struct {
		ctx      context.Context
		storage  Storage
		filePath string

		wantErr error
	}

	tests := map[string]tcase{
		"OK": {
			ctx: context.Background(),
			storage: &mockStorage{
				WritePortFunc: func(ctx context.Context, id string, port *Port) error {
					return nil
				},
			},
			filePath: "testdata/ports.json",
			wantErr:  nil,
		},
		"ErrStorageErr": {
			ctx: context.Background(),
			storage: &mockStorage{
				WritePortFunc: func(ctx context.Context, id string, port *Port) error {
					return errors.New("test error")
				},
			},
			filePath: "testdata/ports.json",
			wantErr:  fmt.Errorf("save port to storage: %w", errors.New("test error")),
		},
		"ErrFileDoesNotExist": {
			ctx: context.Background(),
			storage: &mockStorage{
				WritePortFunc: func(ctx context.Context, id string, port *Port) error {
					return nil
				},
			},
			filePath: "testdata/notexists",
			wantErr: fmt.Errorf("json file streamer: open '%s' file: %w",
				"testdata/notexists", &os.PathError{
					Op:   "open",
					Path: "testdata/notexists",
					Err:  syscall.ENOENT,
				}),
		},
		"ErrInvalid": {
			ctx: context.Background(),
			storage: &mockStorage{
				WritePortFunc: func(ctx context.Context, id string, port *Port) error {
					return nil
				},
			},
			filePath: "testdata/invalid.json",
			wantErr:  fmt.Errorf("decode port id token: %w", io.EOF),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := NewProcessor(tc.storage)

			err := p.Process(tc.ctx, tc.filePath)

			td.Cmp(t, err, tc.wantErr,
				fmt.Sprintf("Process(): expected error := %v, got := %v", tc.wantErr, err),
			)
		})
	}
}

type mockStorage struct {
	WritePortFunc func(ctx context.Context, id string, port *Port) error
}

func (m *mockStorage) WritePort(ctx context.Context, id string, port *Port) error {
	return m.WritePortFunc(ctx, id, port)
}
