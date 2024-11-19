package main

import (
	"context"
	"errors"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

// library is an aggregate root that represents a library.
type library struct {
	books map[string]*book
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
	return l.books[title]
}

// Load loads events with context.
func (l *library) Load(ctx context.Context, events ...aggregate.DomainEventInterface) error {
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
	l.books[event.Title] = &book{title: event.Title}
	return nil
}

// listenDomainEvents listens for domain events.
func (l *library) handleRemoveBookEvent(event *removeBookEvent) error {
	if event.Title == "" {
		return errors.New("title is required")
	}
	delete(l.books, event.Title)
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
	for _, b := range l.books {
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
		books:         make(map[string]*book),
		AggregateRoot: &aggregate.AggregateRoot{},
	}
}
