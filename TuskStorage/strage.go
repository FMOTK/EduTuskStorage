package tuskstorage

import (
	"context"
	"fmt"
	"log"
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

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return

			default:
				s.mu.Lock()
				for id, t := range s.data {
					currTime := time.Now().UTC()

					if currTime.After(t.expireAt) {
						delete(s.data, id)
						log.Printf("Tusk id=\"%s\" deleted", id)
					}
				}
				s.mu.Unlock()

				time.Sleep(2 * time.Second)
			}
		}
	}()

	return s
}

func (s *TuskStorage) CreateTusk(duration time.Duration) *Tusk {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := NewTask(duration)
	tuskID := t.GetUUID()
	s.data[tuskID] = t

	return t
}

func (s *TuskStorage) UpdateTuskById(id string, status TuskStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tusk, ok := s.data[id]; !ok {
		return fmt.Errorf("Tusk with id=\"%s\"not found", id)
	} else {
		tusk.status = status
	}
	return nil
}

func (s *TuskStorage) GetTuskStatuById(id string) (TuskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if tusk, ok := s.data[id]; !ok {
		return "", fmt.Errorf("Tusk with id=\"%s\"not found", id)
	} else {
		return tusk.status, nil
	}
}

func (s *TuskStorage) DeleteTuskById(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)

	return nil
}
