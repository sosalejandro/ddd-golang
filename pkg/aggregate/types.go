//go:generate mockgen -source=types.go -destination=../mocks/aggregate_root_mock.go -package=mocks
package aggregate

import (
	"context"

	"github.com/google/uuid"
)

// AggregateRootInterface is an interface that represents an aggregate root.
type AggregateRootInterface interface {
	Id() uuid.UUID
	SetId(uuid.UUID)
	GetChanges() []RecordedEvent
	Load(ctx context.Context, events ...RecordedEvent) error
	Deserialize([]byte) error
	Serialize() ([]byte, error)
}
