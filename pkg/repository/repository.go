package repository

// import (
// 	"context"
// 	"encoding/json"
// 	"reflect"

// 	"github.com/google/uuid"
// 	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
// 	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
// )

// // RehydrateAggregate is a function that rehydrates an aggregate.
// // type RehydrateAggregate func(context.Context, uuid.UUID, aggregate.AggregateRootInterface) error

// // SnapshotAggregate is a function that saves a snapshot of an aggregate.
// // type SnapshotAggregate func(context.Context, aggregate.AggregateRootInterface) error

// // SaveEvents is a function that saves an aggregate events.
// type SaveEvents func(context.Context, []pkg.EventPayload) error

// // LoadSnapshot is a function that loads a snapshot of an aggregate.
// type LoadSnapshot func(context.Context, uuid.UUID) (*Snapshot, error)

// // GetAggregateEvents is a function that gets aggregate events.
// type GetAggregateEvents func(context.Context, uuid.UUID) ([]pkg.EventPayload, error)

// // SaveSnapshot is a function that saves a snapshot.
// type SaveSnapshot func(context.Context, []byte) error

// // Close is a function that closes the repository.
// type Close func() error

// type RepositoryInterface interface {
// 	// Save saves an aggregate root.
// 	Save(ctx context.Context, aggregateRoot aggregate.AggregateRootInterface) error
// 	// Rehydrate rehydrates an aggregate root.
// 	Rehydrate(ctx context.Context, aggregateId uuid.UUID) (aggregate.AggregateRootInterface, error)
// 	// Snapshot saves a snapshot of an aggregate root.
// 	Snapshot(ctx context.Context, aggregateRoot aggregate.AggregateRootInterface) error
// 	// LoadSnapshot loads a snapshot of an aggregate root.
// 	Load(ctx context.Context, aggregateId uuid.UUID) (aggregate.AggregateRootInterface, error)
// }

// // type RepositoryGeneric[T any] struct {
// // 	CreateFunc func(ctx context.Context, entity T) error
// // 	ReadFunc   func(ctx context.Context, id string) (T, error)
// // 	UpdateFunc func(ctx context.Context, entity T) error
// // 	DeleteFunc func(ctx context.Context, id string) error
// // }

// type Repository struct {
// 	em *pkg.EventManager
// }

// func NewRepository(em *pkg.EventManager) *Repository {
// 	return &Repository{
// 		em: em,
// 	}
// }

// func (r *Repository) Save(ctx context.Context, aggregateRoot aggregate.AggregateRootInterface, cb SaveEvents) error {
// 	changes := aggregateRoot.GetChanges()
// 	if len(changes) == 0 {
// 		return nil
// 	}

// 	var changesToSave []pkg.EventPayload
// 	aggregateId := aggregateRoot.Id()

// 	for _, change := range changes {
// 		data, err := json.Marshal(change.Event)
// 		if err != nil {
// 			return err
// 		}

// 		changesToSave = append(changesToSave, pkg.EventPayload{
// 			AggregateID: aggregateId,
// 			Timestamp:   change.Timestamp,
// 			EventType:   reflect.TypeOf(change.Event).Elem().Name(),
// 			EventID:     change.EventID,
// 			Data:        data,
// 		})
// 	}

// 	if err := cb(ctx, changesToSave); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *Repository) Rehydrate(ctx context.Context, aggregateId uuid.UUID, aggregateRoot aggregate.AggregateRootInterface, cb GetAggregateEvents) error {
// 	events, err := cb(ctx, aggregateId)
// 	if err != nil {
// 		return err
// 	}

// 	domainEvents := make([]aggregate.RecordedEvent, 0)
// 	for _, event := range events {
// 		eventType := event.EventType
// 		data := event.Data

// 		eventInterface, err := r.em.UnmarshalEvent(eventType, data)
// 		if err != nil {
// 			return err
// 		}

// 		domainEvents = append(domainEvents,
// 			aggregate.RecordedEvent{
// 				EventID:   event.EventID,
// 				Timestamp: event.Timestamp,
// 				Event:     eventInterface.(aggregate.DomainEventInterface),
// 			})
// 	}

// 	aggregateRoot.SetId(aggregateId)
// 	if err := aggregateRoot.Load(ctx, domainEvents...); err != nil {
// 		return err
// 	}

// 	return nil
// }
