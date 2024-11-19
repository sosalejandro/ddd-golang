package aggregate_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	"github.com/stretchr/testify/assert"
)

// addBookEvent is a domain event that represents adding a book to the library.
type addBookEvent struct {
	Title string
}

// removeBookEvent is a domain event that represents removing a book from the library.
type removeBookEvent struct {
	Title string
}

// library is an aggregate root that represents a library.
type library struct {
	books map[string]*book
	*aggregate.AggregateRoot
}

// AddBook adds a book to the library with context.
func (l *library) AddBook(ctx context.Context, title string) error {
	return l.ApplyDomainEvent(ctx, &addBookEvent{Title: title}, l.Handle)
}

// RemoveBook removes a book from the library with context.
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
		AggregateRoot: aggregate.NewAggregateRoot(),
	}
}

func TestAggregateRoot(t *testing.T) {
	t.Run("should handle valid book operations", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		// Act
		lib.AddBook(context.Background(), "Clean Code")

		// Assert
		assert.NotNil(t, lib.GetBook("Clean Code"))
		assert.Equal(t, 0, lib.Version)
		assert.Len(t, lib.GetChanges(), 1)
	})

	t.Run("should detect invalid book state", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		// Act & Assert
		err := lib.AddBook(context.Background(), "")

		assert.Error(t, err)
		// Verify state wasn't changed
		assert.Equal(t, -1, lib.Version)
		assert.Empty(t, lib.GetChanges())
		assert.Nil(t, lib.GetBook(""))
	})

	t.Run("should maintain version across operations", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		// Act & Assert
		lib.AddBook(context.Background(), "Book1")
		assert.Equal(t, 0, lib.Version)

		lib.AddBook(context.Background(), "Book2")
		assert.Equal(t, 1, lib.Version)

		lib.RemoveBook(context.Background(), "Book1")
		assert.Equal(t, 2, lib.Version)
		assert.Nil(t, lib.GetBook("Book1"))
	})

	t.Run("should load history correctly", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		history := []aggregate.DomainEventInterface{
			&addBookEvent{Title: "Book1"},
			&addBookEvent{Title: "Book2"},
			&removeBookEvent{Title: "Book1"},
		}

		// Act
		lib.Load(context.Background(), history...)

		// Assert
		assert.Equal(t, 2, lib.Version)
		assert.Nil(t, lib.GetBook("Book1"))
		assert.NotNil(t, lib.GetBook("Book2"))
		assert.Empty(t, lib.GetChanges())
	})

	t.Run("should handle error when loading history", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		history := []aggregate.DomainEventInterface{
			&addBookEvent{Title: "Book1"},
			&addBookEvent{Title: "Book2"},
			&removeBookEvent{Title: "Book1"},
			&removeBookEvent{Title: ""},
		}

		// Act
		err := lib.Load(context.Background(), history...)

		// Assert
		assert.Error(t, err)
	})

	t.Run("should handle sequential add and remove operations", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		expectedEvents := 500

		// Act
		for i := 0; i < 250; i++ {
			lib.AddBook(context.Background(), fmt.Sprintf("Book%d", i))
			lib.RemoveBook(context.Background(), fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.Equal(t, 499, lib.Version)
		assert.Len(t, lib.GetChanges(), expectedEvents)
		for i := 0; i < 250; i++ {
			assert.Nil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential add operations", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		expectedEvents := 500

		// Act
		for i := 0; i < 500; i++ {
			lib.AddBook(context.Background(), fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.Equal(t, 499, lib.Version)
		assert.Len(t, lib.GetChanges(), expectedEvents)
		for i := 0; i < 500; i++ {
			assert.NotNil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential remove operations", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		expectedEvents := 1000 // 500 adds + 500 removes

		// Add books first
		for i := 0; i < 500; i++ {
			lib.AddBook(context.Background(), fmt.Sprintf("Book%d", i))
		}

		// Act
		for i := 0; i < 500; i++ {
			lib.RemoveBook(context.Background(), fmt.Sprintf("Book%d", i))
		}
		// Assert
		assert.Equal(t, 999, lib.Version)
		assert.Len(t, lib.GetChanges(), expectedEvents)
		for i := 0; i < 500; i++ {
			assert.Nil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}

	})

	t.Run("should clear changes without affecting state", func(t *testing.T) {
		// Arrange

		lib := newLibrary()

		// Act
		lib.AddBook(context.Background(), "Test Book")
		initialVersion := lib.Version
		lib.ClearChanges()

		// Assert
		assert.NotNil(t, lib.GetBook("Test Book"))
		assert.Empty(t, lib.GetChanges())
		assert.Equal(t, initialVersion, lib.Version)
	})

	t.Run("should handle valid book operations with context", func(t *testing.T) {
		// Arrange
		lib := newLibrary()

		// Act
		lib.AddBook(context.Background(), "Clean Code")

		// Assert
		assert.Eventually(t, func() bool {
			return lib.GetBook("Clean Code") != nil
		}, 50*time.Millisecond, 5*time.Millisecond)
		assert.Equal(t, 0, lib.Version)
		assert.Len(t, lib.GetChanges(), 1)
	})
}

func TestAggregateRoot_WithContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation during AddBook", func(t *testing.T) {
		// Arrange
		lib := newLibrary()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		// Act
		err := lib.AddBook(ctx, "Clean Code")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, -1, lib.Version)
		assert.Empty(t, lib.GetChanges())
		assert.Nil(t, lib.GetBook("Clean Code"))
	})

	t.Run("should handle context cancellation during Load", func(t *testing.T) {
		// Arrange
		lib := newLibrary()
		ctx, cancel := context.WithCancel(context.Background())
		history := []aggregate.DomainEventInterface{
			&addBookEvent{Title: "Book1"},
			&addBookEvent{Title: "Book2"},
		}
		cancel() // Cancel the context before loading

		// Act
		err := lib.Load(ctx, history...)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, -1, lib.Version)
		assert.Empty(t, lib.GetChanges())
		assert.Nil(t, lib.GetBook("Book1"))
		assert.Nil(t, lib.GetBook("Book2"))
	})
}

func BenchmarkAggregateRoot(b *testing.B) {
	b.Run("Sequential AddBook", func(b *testing.B) {
		lib := newLibrary()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lib.AddBook(context.Background(), "Book")
		}
	})

	b.Run("Load History", func(b *testing.B) {
		lib := newLibrary()

		history := []aggregate.DomainEventInterface{
			&addBookEvent{Title: "Book1"},
			&removeBookEvent{Title: "Book1"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lib.Load(context.Background(), history...)
		}
	})
}

// func BenchmarkParallelAggregateRoot(b *testing.B) {
// 	b.Run("Parallel AddBook", func(b *testing.B) {

// 		lib := newLibrary()

// 		b.ResetTimer()
// 		b.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				lib.AddBook("Book")
// 			}
// 		})
// 	})
// }
