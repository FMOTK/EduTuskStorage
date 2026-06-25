package tuskmanager

import (
	"context"
	"fmt"

	storage "httpServer/TuskStorage"

	"github.com/google/uuid"
)

type TuskManager struct {
	storage *storage.TuskStorage
	ctx     context.Context
}

func NewManager(ctx context.Context, storage *storage.TuskStorage) *TuskManager {
	return &TuskManager{
		storage: storage,
		ctx:     ctx,
	}
}

func (m *TuskManager) StartTusk(t *storage.Tusk) error {
	//Validate Tusk
	id := t.GetUUID()
	if err := uuid.Validate(id); err != nil {
		return fmt.Errorf("Tusk start error: %v", err)
	}
	if m.storage.IsContain(id) {
		return fmt.Errorf("Tusk UUID is not uniq for storage")
	}
	if t == nil {
		return fmt.Errorf("Tusk is nil")
	}

	m.storage.Set(t)
	go t.Run(m.ctx)

	return nil
}
