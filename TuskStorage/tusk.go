package tuskstorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TuskStatus string

const (
	RunningStatus   TuskStatus = "running"
	CompletedStatus TuskStatus = "complited"
	FailedStatus    TuskStatus = "failed"
	PendingStatus   TuskStatus = "pending"
	CancelStatus    TuskStatus = "cancel"
)

var (
	ErrClientCanceled = errors.New("client canceled")
	ErrServerStoped   = errors.New("stoped canceled")
)

type Tusk struct {
	uuid     string
	status   TuskStatus
	duration time.Duration

	createdAt time.Time
	expireAt  time.Time

	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelCauseFunc
}

func NewTask(duration time.Duration, TTL string) (*Tusk, error) {

	if duration < 0 {
		return nil, fmt.Errorf("duration < 0")
	}

	ctx, cancel := context.WithCancelCause(context.Background())

	id := uuid.New().String()

	currTime := time.Now().UTC()
	ttl, _ := time.ParseDuration(TTL)
	deadTime := currTime.Add(ttl)

	return &Tusk{
		uuid:      id,
		duration:  duration,
		status:    PendingStatus,
		createdAt: currTime,
		expireAt:  deadTime,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (t *Tusk) GetExpiredTime() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.expireAt
}

func (t *Tusk) GetUUID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.uuid
}

func (t *Tusk) GetStatus() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return string(t.status)
}

func (t *Tusk) setStatus(status TuskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status = status
}

func (t *Tusk) Run(ctx context.Context) {

	t.mu.Lock()
	t.ctx, t.cancel = context.WithCancelCause(ctx)
	duration := t.duration
	t.mu.Unlock()

	t.setStatus(RunningStatus)

	select {
	case <-t.ctx.Done():
		cause := context.Cause(t.ctx)
		if cause == ErrClientCanceled {
			t.setStatus(CancelStatus)
		} else {
			t.setStatus(FailedStatus)
		}
	case <-time.After(duration):
		t.setStatus(CompletedStatus)
		t.cancel(nil)
	}
}

// TODO: Отменя задачи через контекст
func (t *Tusk) Cancel() error {
	select {
	case <-t.ctx.Done():
		return nil
	default:
		t.cancel(ErrClientCanceled)
		return nil
	}
}
