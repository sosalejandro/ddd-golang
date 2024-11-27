package repository

import (
	"context"
	"encoding/json"
	"reflect"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	event_manager "github.com/sosalejandro/ddd-golang/pkg/event-manager"
)

// GenericRepository is a generic repository for aggregates.
type GenericRepository[T aggregate.AggregateRootInterface] struct {
	em      *event_manager.EventManager
	repo    AggregateRepositoryInterface
	factory func() T
}

func NewGenericRepository[T aggregate.AggregateRootInterface](em *event_manager.EventManager, funcs AggregateRepositoryInterface, factory func() T) *GenericRepository[T] {
	return &GenericRepository[T]{
		em:      em,
		repo:    funcs,
		factory: factory,
	}
}

func (r *GenericRepository[T]) SaveEvents(ctx context.Context, aggregateRoot T) error {
	changes := aggregateRoot.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	var changesToSave []event_manager.EventPayload
	aggregateId := aggregateRoot.Id()

	for _, change := range changes {
		data, err := json.Marshal(change.Event)
		if err != nil {
			return err
		}

		changesToSave = append(changesToSave, event_manager.EventPayload{
			AggregateID: aggregateId,
			Timestamp:   change.Timestamp,
			EventType:   reflect.TypeOf(change.Event).Elem().Name(),
			EventID:     change.EventID,
			Data:        data,
		})
	}

	return r.repo.SaveEvents(ctx, changesToSave)
}

func IterateIntoRecordedEvents(events []event_manager.EventPayload, eventManager *event_manager.EventManager, iterErr *error) func(func(aggregate.RecordedEvent) bool) {
	return func(yield func(aggregate.RecordedEvent) bool) {
		for _, event := range events {
			eventInterface, err := eventManager.UnmarshalEvent(event.EventType, event.Data)
			if err != nil {
				*iterErr = err
				return
			}

			if !yield(aggregate.RecordedEvent{
				EventID:   event.EventID,
				Timestamp: event.Timestamp,
				Event:     eventInterface.(aggregate.DomainEventInterface),
			}) {
				return
			}
		}
	}
}

func (r *GenericRepository[T]) Rehydrate(ctx context.Context, aggregateId uuid.UUID) (T, error) {
	aggregateRoot := r.factory()

	events, err := r.repo.GetAggregateEvents(ctx, aggregateId)
	if err != nil {
		return aggregateRoot, err
	}

	// domainEvents := make([]aggregate.RecordedEvent, 0)
	iterErr := new(error)
	domainEvents := slices.Collect(IterateIntoRecordedEvents(events, r.em, iterErr))

	if *iterErr != nil {
		return aggregateRoot, *iterErr
	}

	aggregateRoot.SetId(aggregateId)
	if err := aggregateRoot.Load(ctx, domainEvents...); err != nil {
		return aggregateRoot, err
	}

	return aggregateRoot, nil
}

func (r *GenericRepository[T]) Load(ctx context.Context, aggregateId uuid.UUID) (T, error) {
	aggregateRoot := r.factory()

	// Load snapshot
	snapshot, err := r.repo.LoadSnapshot(ctx, aggregateId)
	if err != nil {
		return aggregateRoot, err
	}

	// Set the ID and deserialize snapshot
	aggregateRoot.SetId(aggregateId) // Assuming SetId is a method of T
	if err := aggregateRoot.Deserialize(snapshot.Data); err != nil {
		return aggregateRoot, err
	}

	return aggregateRoot, nil
}

func (r *GenericRepository[T]) SaveSnapshot(ctx context.Context, aggregateRoot T) error {
	data, err := aggregateRoot.Serialize()
	if err != nil {
		return err
	}

	getLatestVersion := func() int {
		return aggregateRoot.GetChanges()[len(aggregateRoot.GetChanges())-1].EventID
	}

	snapshot := &Snapshot{
		AggregateID: aggregateRoot.Id(),
		Data:        data,
		Version:     getLatestVersion(),
		Timestamp:   time.Now().Unix(),
	}

	return r.repo.SaveSnapshot(ctx, snapshot)
}

// Close manually disposes the internal connection
func (r *GenericRepository[T]) Close() error {
	if r.repo != nil {
		return r.repo.Close()
	}
	return nil
}
