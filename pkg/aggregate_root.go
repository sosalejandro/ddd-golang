package pkg

import (
	"context"
	"sync/atomic"
)

type IAggregateRoot interface {
	Handle(event IDomainEvent)
	ValidateState() error
	SetDomainEventsChannel(ch chan<- IDomainEvent)
	SendDomainEvent(event IDomainEvent)
}

type AggregateRoot[T IAggregateRoot] struct {
	changes         []IDomainEvent
	Version         int
	domainEventCh   chan IDomainEvent
	entity          T
	errCh           chan error
	processedEvents int32 // atomic counter
}

func NewAggregateRoot[T IAggregateRoot](ctx context.Context, entity T, errCh chan error) (*AggregateRoot[IAggregateRoot], func()) {
	ar := &AggregateRoot[IAggregateRoot]{
		changes:       make([]IDomainEvent, 0),
		Version:       -1,
		entity:        entity,
		domainEventCh: make(chan IDomainEvent),
		errCh:         errCh,
	}
	ar.entity.SetDomainEventsChannel(ar.domainEventCh)
	go ar.listenDomainEvents(ctx)
	return ar, ar.closeChannel
}

func (ar *AggregateRoot[T]) GetProcessedEventsCount() int {
	return int(atomic.LoadInt32(&ar.processedEvents))
}

func (ar *AggregateRoot[T]) closeChannel() {
	close(ar.domainEventCh)
}

func (ar *AggregateRoot[T]) GetChanges() []IDomainEvent {
	return ar.changes
}

func (ar *AggregateRoot[T]) ClearChanges() {
	ar.changes = make([]IDomainEvent, 0)
}

func (ar *AggregateRoot[T]) ApplyDomainEvent(event IDomainEvent) error {
	ar.entity.Handle(event)
	if err := ar.entity.ValidateState(); err != nil {
		return err
	}
	ar.changes = append(ar.changes, event)
	ar.Version++
	atomic.AddInt32(&ar.processedEvents, 1)
	return nil
}

func (ar *AggregateRoot[T]) Load(history []IDomainEvent) {
	for _, e := range history {
		ar.entity.Handle(e)
		ar.Version++
	}
	ar.ClearChanges()
}

func (ar *AggregateRoot[T]) listenDomainEvents(ctx context.Context) {
	for {
		select {
		case event := <-ar.domainEventCh:
			if err := ar.ApplyDomainEvent(event); err != nil {
				ar.errCh <- err
			}
		case <-ctx.Done():
			ar.closeChannel()
			return
		}
	}
}
