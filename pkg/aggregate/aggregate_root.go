package aggregate

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// AggregateRootInterface is an interface that represents an aggregate root entity.
type AggregateRootInterface interface {
	// Handle handles domain events.
	Handle(event IDomainEvent) error
	// ValidateState validates the state of the aggregate root.
	ValidateState() error
	// SetDomainEventsChannel sets the channel to send domain events.
	SetDomainEventsChannel(ch chan<- IDomainEvent)
	// SendDomainEvent sends a domain event.
	SendDomainEvent(event IDomainEvent)
}

// AggregateRoot is a struct that represents an aggregate root.
// It keeps track of the changes, the version, the domain events channel, the entity, the error channel, and the processed events.
// It also has methods to get the processed events count, close the domain events channel, get the changes, clear the changes, apply a domain event, load the history, and listen for domain events.
type AggregateRoot[T AggregateRootInterface] struct {
	changes          []IDomainEvent    // Recorded domain events
	Version          int               // Current version of the aggregate
	domainEventCh    chan IDomainEvent // Channel for incoming domain events
	entity           T                 // The entity managed by the aggregate
	errCh            chan error        // Channel for reporting errors
	processedEvents  int32             // Atomic counter for processed events
	closeOnce        sync.Once         // Ensures channels are closed only once
	eventProcessedCh chan struct{}     // Signals when an event has been processed
}

// NewAggregateRoot creates a new aggregate root.
// It returns the aggregate root and a function to close the domain events channel.
// The aggregate root listens for domain events and applies them.
// The aggregate root also keeps track of the changes.
// The aggregate root also keeps track of the version.
// The aggregate root also keeps track of the processed events.
func NewAggregateRoot[T AggregateRootInterface](ctx context.Context, entity T, errCh chan error) (*AggregateRoot[AggregateRootInterface], func()) {
	ar := &AggregateRoot[AggregateRootInterface]{
		changes:          make([]IDomainEvent, 0),
		Version:          -1,
		entity:           entity,
		domainEventCh:    make(chan IDomainEvent, 100), // buffered
		errCh:            errCh,                        // buffered
		eventProcessedCh: make(chan struct{}, 1),       // buffered channel to prevent blocking
	}
	ar.entity.SetDomainEventsChannel(ar.domainEventCh)
	go ar.listenDomainEvents(ctx)
	return ar, ar.closeChannel
}

// GetProcessedEventsCount returns the number of processed events.
func (ar *AggregateRoot[T]) GetProcessedEventsCount() int {
	return int(atomic.LoadInt32(&ar.processedEvents))
}

// closeChannel closes the channels.
// It closes the domain events channel, the error channel, and the event processed channel.
// It closes the channels only once.
func (ar *AggregateRoot[T]) closeChannel() {
	ar.closeOnce.Do(func() {
		close(ar.domainEventCh)
		close(ar.errCh)
		close(ar.eventProcessedCh)
	})
}

// GetChanges returns the changes.
func (ar *AggregateRoot[T]) GetChanges() []IDomainEvent {
	return ar.changes
}

// ClearChanges clears the changes.
// It sets the changes to an empty slice.
func (ar *AggregateRoot[T]) ClearChanges() {
	ar.changes = make([]IDomainEvent, 0)
}

// ApplyDomainEvent applies a domain event.
// It validates the state of the entity.
// It increments the version.
// It increments the processed events.
// It returns an error if the state is invalid.
func (ar *AggregateRoot[T]) ApplyDomainEvent(event IDomainEvent) (err error) {
	if err = ar.entity.Handle(event); err != nil {
		return err
	}

	if err = ar.entity.ValidateState(); err != nil {
		return err
	}

	// Signal that an event was processed, regardless of validation result
	atomic.AddInt32(&ar.processedEvents, 1)
	select {
	case ar.eventProcessedCh <- struct{}{}:
	default:
	}

	ar.changes = append(ar.changes, event)
	ar.Version++
	return nil
}

// Load loads the history.
// It applies the domain events to the entity.
// It increments the version.
// It clears the changes.
func (ar *AggregateRoot[T]) Load(history []IDomainEvent) {
	for _, e := range history {
		if err := ar.ApplyDomainEvent(e); err != nil {
			ar.errCh <- err
		}
	}
	ar.ClearChanges()
}

// listenDomainEvents listens for domain events.
// It listens for domain events and applies them.
// It closes the domain events channel when the context is done.
// It sends an error to the error channel if there is an error.
// listenDomainEvents listens for incoming domain events and processes them.
func (ar *AggregateRoot[T]) listenDomainEvents(ctx context.Context) {
	for {
		select {
		case event, ok := <-ar.domainEventCh:
			if !ok {
				ar.closeChannel()
				return
			}
			if err := ar.ApplyDomainEvent(event); err != nil {
				ar.errCh <- err
			}
		case <-ctx.Done():
			ar.closeChannel()
			return
		}
	}
}

// WaitForEvents waits until the target count of events is processed or the context is canceled.
func (ar *AggregateRoot[T]) WaitForEvents(ctx context.Context, targetCount int) error {
	for ar.GetProcessedEventsCount() < targetCount {
		select {
		case <-ar.eventProcessedCh:
			// Continue waiting
		case <-ctx.Done():
			return fmt.Errorf("operation canceled while waiting for events")
		}
	}
	return nil
}
