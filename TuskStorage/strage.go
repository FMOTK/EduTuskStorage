package tuskstorage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TuskStorage struct {
	mu   sync.RWMutex
	data map[string]*Tusk
	ctx  context.Context
}

func NewTuskStorage(ctx context.Context) *TuskStorage {
	s := &TuskStorage{
		data: make(map[string]*Tusk),
		ctx:  ctx,
	}
	return s
}

func (s *TuskStorage) GetExpireds() map[string]*Tusk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]*Tusk)
	currTime := time.Now().UTC()

	for _, t := range s.data {
		if currTime.After(t.expireAt) {
			out[t.GetUUID()] = t
		}
	}

	return out
}

func (s *TuskStorage) Set(t *Tusk) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tuskID := t.GetUUID()
	s.data[tuskID] = t
}

func (s *TuskStorage) Get(id string) (*Tusk, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if tusk, ok := s.data[id]; !ok {
		return nil, fmt.Errorf("Tusk with id=\"%s\"not found", id)
	} else {
		result := tusk
		return result, nil
	}
}

func (s *TuskStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tusk, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no tusk ID=%s in storage", id)
	}

	tusk.Cancel()
	delete(s.data, id)

	return nil
}

func (s *TuskStorage) IsContain(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[id]
	return ok
}
