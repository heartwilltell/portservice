package inmemstore

import (
	"context"
	"testing"

	"jsonstream/app/port"

	"github.com/maxatome/go-testdeep/td"
)

func TestStorage_WritePort(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		p := &port.Port{Name: "test port"}

		s := New()

		err := s.WritePort(context.Background(), "test", p)

		td.CmpNoError(t, err)
		td.Cmp(t, s.store, map[string]*port.Port{
			"test": p,
		})
	})
}
