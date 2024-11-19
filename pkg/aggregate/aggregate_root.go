package aggregate

import (
	"context"
)

// AggregateRootInterface is an interface that represents an aggregate root.
type AggregateRootInterface interface {
	// Handle handles a domain event with context and returns an error if there is an error.
	Handle(ctx context.Context, event DomainEventInterface) error
	// Load loads the history of domain events with context to the entity and applies them to the entity.
	Load(ctx context.Context, history ...DomainEventInterface) error
}

// AggregateRoot is a struct that represents an aggregate root.
// It keeps track of the changes.
// It keeps track of the version.
type AggregateRoot struct {
	changes []DomainEventInterface // Recorded domain events
	Version int                    // Current version of the aggregate
}

// NewAggregateRoot creates a new aggregate root.
// It returns the aggregate root and a function to close the domain events channel.
// The aggregate root also keeps track of the changes.
// The aggregate root also keeps track of the version.
func NewAggregateRoot() *AggregateRoot {
	ar := &AggregateRoot{
		changes: make([]DomainEventInterface, 0),
		Version: -1,
	}
	return ar
}

// GetChanges returns the changes.
func (ar *AggregateRoot) GetChanges() []DomainEventInterface {
	return ar.changes
}

// ClearChanges clears the changes.
// It sets the changes to an empty slice.
func (ar *AggregateRoot) ClearChanges() {
	ar.changes = make([]DomainEventInterface, 0)
}

// ApplyDomainEvent applies a domain event with context.
// It checks for context cancellation before handling the event.
func (ar *AggregateRoot) ApplyDomainEvent(ctx context.Context, event DomainEventInterface, handle func(ctx context.Context, e DomainEventInterface) error) (err error) {
	// Check if the context is already canceled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err = handle(ctx, event); err != nil {
		return err
	}

	ar.changes = append(ar.changes, event)
	ar.Version++
	return nil
}

// Load loads the history with context.
// It checks for context cancellation during the loading process.
func (ar *AggregateRoot) Load(ctx context.Context, history []DomainEventInterface, handle func(ctx context.Context, e DomainEventInterface) error) error {
	for _, e := range history {
		// Check for context cancellation before applying each event
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := ar.ApplyDomainEvent(ctx, e, handle); err != nil {
			return err
		}
	}
	ar.ClearChanges()
	return nil
}
