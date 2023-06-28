package inmemstore

import (
	"context"
	"sync"

	"jsonstream/app/port"
)

// Storage implements port.Storage interface using simple in-memory hash table.
// Please don't use such storage for production because of its inefficiency.
type Storage struct {
	mu    sync.RWMutex
	store map[string]*port.Port
}

// New returns a pointer to a new instance of Storage.
func New() *Storage { return &Storage{store: make(map[string]*port.Port)} }

func (s *Storage) WritePort(_ context.Context, id string, port *port.Port) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[id] = port
	return nil
}
