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
	ch    chan<- aggregate.IDomainEvent
}

// AddBook adds a book to the library.
func (l *library) AddBook(title string) {
	l.SendDomainEvent(&addBookEvent{Title: title})
}

// RemoveBook removes a book from the library.
func (l *library) RemoveBook(title string) {
	l.SendDomainEvent(&removeBookEvent{Title: title})
}

// GetBook returns a book from the library.
func (l *library) GetBook(title string) *book {
	return l.books[title]
}

// SetDomainEventsChannel sets the channel to send domain events.
func (l *library) SetDomainEventsChannel(ch chan<- aggregate.IDomainEvent) {
	l.ch = ch
}

// SendDomainEvent sends a domain event.
func (l *library) SendDomainEvent(event aggregate.IDomainEvent) {
	l.ch <- event
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
func (l *library) Handle(event aggregate.IDomainEvent) error {
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
		books: make(map[string]*book),
	}
}

func TestAggregateRoot(t *testing.T) {
	t.Run("should handle valid book operations", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Clean Code")
		assert.NoError(t, ar.WaitForEvents(ctx, 1))

		// Assert
		assert.NotNil(t, lib.GetBook("Clean Code"))
		assert.Equal(t, 0, ar.Version)
		assert.Len(t, ar.GetChanges(), 1)
	})

	t.Run("should detect invalid book state", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		lib := newLibrary()
		errCh := make(chan error, 10)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		// Act & Assert
		lib.AddBook("")

		// Wait for event processing
		assert.NoError(t, ar.WaitForEvents(ctx, 0))

		// Verify error was received
		select {
		case err := <-errCh:
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "title is required")
		case <-time.After(time.Second):
			t.Fatal("expected error for invalid book")
		}

		// Verify state wasn't changed
		assert.Equal(t, -1, ar.Version)
		assert.Empty(t, ar.GetChanges())
		assert.Nil(t, lib.GetBook(""))
	})

	t.Run("should maintain version across operations", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		// Act & Assert
		lib.AddBook("Book1")
		assert.NoError(t, ar.WaitForEvents(ctx, 1))
		assert.Equal(t, 0, ar.Version)

		lib.AddBook("Book2")
		assert.NoError(t, ar.WaitForEvents(ctx, 2))
		assert.Equal(t, 1, ar.Version)

		lib.RemoveBook("Book1")
		assert.NoError(t, ar.WaitForEvents(ctx, 3))
		assert.Equal(t, 2, ar.Version)
	})

	t.Run("should load history correctly", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		history := []aggregate.IDomainEvent{
			&addBookEvent{Title: "Book1"},
			&addBookEvent{Title: "Book2"},
			&removeBookEvent{Title: "Book1"},
		}

		// Act
		ar.Load(history)

		// Assert
		assert.Equal(t, 2, ar.Version)
		assert.Nil(t, lib.GetBook("Book1"))
		assert.NotNil(t, lib.GetBook("Book2"))
		assert.Empty(t, ar.GetChanges())
		assert.Equal(t, 3, ar.GetProcessedEventsCount())
	})

	t.Run("should handle sequential add and remove operations", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		expectedEvents := 500

		// Act
		for i := 0; i < 250; i++ {
			lib.AddBook(fmt.Sprintf("Book%d", i))
			lib.RemoveBook(fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.NoError(t, ar.WaitForEvents(ctx, expectedEvents))
		assert.Equal(t, 499, ar.Version)
		assert.Len(t, ar.GetChanges(), expectedEvents)
		for i := 0; i < 250; i++ {
			assert.Nil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential add operations", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		expectedEvents := 500

		// Act
		for i := 0; i < 500; i++ {
			lib.AddBook(fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.NoError(t, ar.WaitForEvents(ctx, expectedEvents))
		assert.Equal(t, 499, ar.Version)
		assert.Len(t, ar.GetChanges(), expectedEvents)
		for i := 0; i < 500; i++ {
			assert.NotNil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential remove operations", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		expectedEvents := 1000 // 500 adds + 500 removes

		// Add books first
		for i := 0; i < 500; i++ {
			lib.AddBook(fmt.Sprintf("Book%d", i))
		}

		// Act
		for i := 0; i < 500; i++ {
			lib.RemoveBook(fmt.Sprintf("Book%d", i))
		}
		// Assert
		assert.NoError(t, ar.WaitForEvents(ctx, expectedEvents))
		assert.Equal(t, 999, ar.Version)
		assert.Len(t, ar.GetChanges(), expectedEvents)
		for i := 0; i < 500; i++ {
			assert.Nil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}

		// Check error channel
		select {
		case err := <-errCh:
			t.Fatalf("unexpected error: %v", err)
		default:
			// No errors - good
		}
	})

	t.Run("should clear changes without affecting state", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Test Book")
		assert.NoError(t, ar.WaitForEvents(ctx, 1))
		initialVersion := ar.Version
		ar.ClearChanges()

		// Assert
		assert.NotNil(t, lib.GetBook("Test Book"))
		assert.Empty(t, ar.GetChanges())
		assert.Equal(t, initialVersion, ar.Version)
	})

	t.Run("should handle valid book operations with context", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Clean Code")
		assert.NoError(t, ar.WaitForEvents(ctx, 1))

		// Assert
		assert.Eventually(t, func() bool {
			return lib.GetBook("Clean Code") != nil
		}, 50*time.Millisecond, 5*time.Millisecond)
		assert.Equal(t, 0, ar.Version)
		assert.Len(t, ar.GetChanges(), 1)
	})

	t.Run("should close channel on context cancellation", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		lib := newLibrary()
		errCh := make(chan error, 10)
		ar, _ := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)

		// Act
		lib.AddBook("Test Book")
		err := ar.WaitForEvents(ctx, 1)
		cancel() // Cancel context

		// Assert - try to send after context cancelled
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond) // Allow time for cleanup

		// This should panic if channel isn't closed - recover and assert
		recovered := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					recovered = true
				}
			}()
			lib.AddBook("Should Fail")
		}()

		assert.True(t, recovered, "expected panic when writing to closed channel")
	})

	t.Run("should timeout if operation takes too long", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		// defer closer()

		// Act & Assert
		time.Sleep(20 * time.Millisecond) // Wait for timeout
		recovered := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					recovered = true
				}
			}()
			lib.AddBook("Should Fail")
		}()

		assert.True(t, recovered, "expected panic when writing to closed channel after timeout")
	})

	t.Run("WaitForEvents should timeout when events are not processed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		err := ar.WaitForEvents(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operation canceled")
	})

	t.Run("WaitForEvents should return immediately when target count is reached", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		defer closer()

		lib.AddBook("Test Book")
		start := time.Now()
		err := ar.WaitForEvents(ctx, 1)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Less(t, duration, 100*time.Millisecond)
	})
}

func BenchmarkAggregateRoot(b *testing.B) {
	b.Run("Sequential AddBook", func(b *testing.B) {
		ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, _ := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		// defer closer()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lib.AddBook("Book")
		}

		if err := ar.WaitForEvents(ctx, b.N); err != nil {
			b.Fatal(err)
		}
	})

	b.Run("Load History", func(b *testing.B) {
		ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, _ := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		// defer closer()

		history := []aggregate.IDomainEvent{
			&addBookEvent{Title: "Book1"},
			&removeBookEvent{Title: "Book1"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ar.Load(history)
		}

		if err := ar.WaitForEvents(ctx, b.N); err != nil {
			b.Fatal(err)
		}
	})
}

func BenchmarkParallelAggregateRoot(b *testing.B) {
	b.Run("Parallel AddBook", func(b *testing.B) {
		ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, _ := aggregate.NewAggregateRoot[aggregate.AggregateRootInterface](ctx, lib, errCh)
		// defer closer()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				lib.AddBook("Book")
			}
		})

		if err := ar.WaitForEvents(ctx, b.N); err != nil {
			b.Fatal(err)
		}
	})
}
