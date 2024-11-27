package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

// library is an aggregate root that represents a library.
type library struct {
	id    uuid.UUID
	Books map[string]*book
	*aggregate.AggregateRoot
}

// AddBook adds a book to the library.
func (l *library) AddBook(ctx context.Context, title string) error {
	return l.ApplyDomainEvent(ctx, &addBookEvent{Title: title}, l.Handle)
}

// RemoveBook removes a book from the library.
func (l *library) RemoveBook(ctx context.Context, title string) error {
	return l.ApplyDomainEvent(ctx, &removeBookEvent{Title: title}, l.Handle)
}

// GetBook returns a book from the library.
func (l *library) GetBook(title string) *book {
	return l.Books[title]
}

// Load loads events with context.
func (l *library) Load(ctx context.Context, events ...aggregate.RecordedEvent) error {
	if l.Books == nil {
		l.Books = make(map[string]*book)
	}

	if l.AggregateRoot == nil {
		l.AggregateRoot = &aggregate.AggregateRoot{}
	}

	if err := l.AggregateRoot.Load(ctx, events, l.Handle); err != nil {
		return err
	}

	l.ClearChanges()
	return nil
}

// listenDomainEvents listens for domain events.
func (l *library) handleAddBookEvent(event *addBookEvent) error {
	if event.Title == "" {
		return errors.New("title is required")
	}
	l.Books[event.Title] = &book{title: event.Title}
	return nil
}

// listenDomainEvents listens for domain events.
func (l *library) handleRemoveBookEvent(event *removeBookEvent) error {
	if event.Title == "" {
		return errors.New("title is required")
	}
	delete(l.Books, event.Title)
	return nil
}

// Handle handles domain events.
func (l *library) Handle(ctx context.Context, event aggregate.DomainEventInterface) error {
	switch e := event.(type) {
	case *addBookEvent:
		return l.handleAddBookEvent(e)
	case *removeBookEvent:
		return l.handleRemoveBookEvent(e)
	}
	return nil
}

// ValidateState validates the state of the library.
func (l *library) ValidateState() error {
	var errs error
	for _, b := range l.Books {
		if err := b.validate(); err != nil {
			if errs == nil {
				errs = err
				continue
			}
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// book is an entity that represents a book.
type book struct {
	title string
}

// validate validates the book.
func (b *book) validate() error {
	if b.title == "" {
		return errors.New("title is required")
	}
	return nil
}

// newLibrary creates a new library.
func newLibrary() *library {
	return &library{
		Books:         make(map[string]*book),
		AggregateRoot: &aggregate.AggregateRoot{},
	}
}

func (l *library) Id() uuid.UUID {
	return l.id
}

func (l *library) SetId(id uuid.UUID) {
	l.id = id
}

func (l *library) Deserialize(data []byte) error {
	if err := json.Unmarshal(data, &l); err != nil {
		return err
	}
	return nil
}

func (l *library) Serialize() ([]byte, error) {
	return json.Marshal(l)
}

// addInvalidBookEvent is a mock event for testing deserialize errors
type addInvalidBookEvent struct {
	// Missing required fields to simulate invalid data
}

// Implement DomainEventInterface
func (e *addInvalidBookEvent) SomeMethod() {}
