package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

type MemoryRepository struct {
	mu        sync.RWMutex
	events    map[uuid.UUID][]pkg.EventPayload
	snapshots map[uuid.UUID]*repository.Snapshot
}

// Ensure MemoryRepository implements AggregateRepositoryInterface
var _ repository.AggregateRepositoryInterface = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		events:    make(map[uuid.UUID][]pkg.EventPayload),
		snapshots: make(map[uuid.UUID]*repository.Snapshot),
	}
}

func (r *MemoryRepository) SaveEvents(ctx context.Context, events []pkg.EventPayload) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, event := range events {
		r.events[event.AggregateID] = append(r.events[event.AggregateID], event)
	}

	return nil
}

func (r *MemoryRepository) GetAggregateEvents(ctx context.Context, id uuid.UUID) ([]pkg.EventPayload, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[id]
	if !ok {
		return nil, fmt.Errorf("events not found for aggregate %s", id)
	}

	return events, nil
}

func (r *MemoryRepository) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID) (*repository.Snapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot, ok := r.snapshots[aggregateID]
	if !ok {
		return nil, fmt.Errorf("snapshot not found for aggregate %s", aggregateID)
	}

	return snapshot, nil
}

func (r *MemoryRepository) SaveSnapshot(ctx context.Context, snapshot *repository.Snapshot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.snapshots[snapshot.AggregateID] = snapshot
	return nil
}

func (r *MemoryRepository) Close() error {
	// No resources to close in memory implementation
	return nil
}
