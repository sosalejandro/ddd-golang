package aggregate

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RecordedEvent wraps a domain event with an EventID representing the aggregate's version
type RecordedEvent struct {
	EventID   uuid.UUID
	Version   int
	Timestamp int64
	Event     DomainEventInterface
}

// AggregateRoot is a struct that represents an aggregate root.
// It keeps track of the changes.
// It keeps track of the version.
type AggregateRoot struct {
	changes []RecordedEvent // Recorded domain events with EventID
	Version int             // Current version of the aggregate
}

// NewAggregateRoot creates a new aggregate root.
// It returns the aggregate root and a function to close the domain events channel.
// The aggregate root also keeps track of the changes.
// The aggregate root also keeps track of the version.
func NewAggregateRoot() *AggregateRoot {
	ar := &AggregateRoot{
		changes: make([]RecordedEvent, 0),
		Version: -1,
	}
	return ar
}

// GetChanges returns the recorded events with EventID.
func (ar *AggregateRoot) GetChanges() []RecordedEvent {
	return ar.changes
}

// ClearChanges clears the recorded events.
// It sets the changes to an empty slice.
func (ar *AggregateRoot) ClearChanges() {
	ar.changes = make([]RecordedEvent, 0)
}

// ApplyDomainEvent applies a domain event with context and assigns an EventID based on the current version.
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

	// Assign EventID as the next version
	eventID := uuid.New()

	// Wrap the event with EventID
	recordedEvent := RecordedEvent{
		EventID:   eventID,
		Timestamp: time.Now().Unix(),
		Event:     event,
		Version:   ar.Version + 1,
	}

	ar.changes = append(ar.changes, recordedEvent)
	ar.Version++
	return nil
}

// Load loads the history with context and assigns EventIDs.
// It checks for context cancellation during the loading process.
func (ar *AggregateRoot) Load(ctx context.Context, history []RecordedEvent, handle func(ctx context.Context, e DomainEventInterface) error) error {
	for _, e := range history {
		// Check for context cancellation before applying each event
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := ar.ApplyDomainEvent(ctx, e.Event, handle); err != nil {
			return err
		}
	}
	ar.ClearChanges()
	return nil
}

// generateUniqueEventID generates a unique identifier for an event
func generateUniqueEventID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
