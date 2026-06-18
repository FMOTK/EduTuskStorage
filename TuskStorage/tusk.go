package tuskstorage

import (
	"context"
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

type Tusk struct {
	uuid     string
	status   TuskStatus
	duration time.Duration

	createdAt time.Time
	expireAt  time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTask(duration time.Duration, TTL string) *Tusk {

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
	}
}

func (t *Tusk) GetUUID() string {
	return t.uuid
}

func (t *Tusk) GetStatus() string {
	return string(t.status)
}

// TODO: Отменя задачи через контекст
func (t *Tusk) Cancel() error {
	select {
	case <-t.ctx.Done():
		return nil
	default:
		t.cancel()
		return nil
	}
}
