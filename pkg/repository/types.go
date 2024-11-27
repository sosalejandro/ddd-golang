//go:generate mockgen -source=types.go -destination=../mocks/aggregate_repository_mock.go -package=mocks
package repository

import (
	"context"

	"github.com/google/uuid"
	event_manager "github.com/sosalejandro/ddd-golang/pkg/event-manager"
)

type AggregateRepositoryInterface interface {
	SaveEvents(ctx context.Context, events []event_manager.EventPayload) error
	GetAggregateEvents(ctx context.Context, id uuid.UUID) ([]event_manager.EventPayload, error)
	LoadSnapshot(ctx context.Context, aggregateID uuid.UUID) (*Snapshot, error)
	SaveSnapshot(ctx context.Context, snapshot *Snapshot) error
	Close() error
}

type Snapshot struct {
	AggregateID uuid.UUID `json:"aggregateId"`
	Version     int       `json:"version"`
	Data        []byte    `json:"data"`
	Timestamp   int64     `json:"timestamp"`
}
