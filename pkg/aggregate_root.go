package pkg

type IAggregateRoot interface {
	Handle(event IDomainEvent)
	ValidateState() error
	SetDomainEventsChannel(ch chan<- IDomainEvent)
	SendDomainEvent(event IDomainEvent)
}

type AggregateRoot[T IAggregateRoot] struct {
	changes       []IDomainEvent
	Version       int
	domainEventCh chan IDomainEvent
	entity        T
	errCh         chan error
}

func NewAggregateRoot[T IAggregateRoot](entity T, errCh chan error) (*AggregateRoot[IAggregateRoot], func()) {
	ar := &AggregateRoot[IAggregateRoot]{
		changes:       make([]IDomainEvent, 0),
		Version:       -1,
		entity:        entity,
		domainEventCh: make(chan IDomainEvent),
		errCh:         errCh,
	}
	ar.entity.SetDomainEventsChannel(ar.domainEventCh)
	go ar.listenDomainEvents()
	return ar, ar.closeChannel
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
	return nil
}

func (ar *AggregateRoot[T]) Load(history []IDomainEvent) {
	for _, e := range history {
		ar.entity.Handle(e)
		ar.Version++
	}
	ar.ClearChanges()
}

func (ar *AggregateRoot[T]) listenDomainEvents() {
	for event := range ar.domainEventCh {
		if err := ar.ApplyDomainEvent(event); err != nil {
			ar.errCh <- err
		}
	}
}
