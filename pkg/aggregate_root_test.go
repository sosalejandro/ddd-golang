package pkg_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	pkg "github.com/sosalejandro/ddd-golang/pkg"
	"github.com/stretchr/testify/assert"
)

type addBookEvent struct {
	Title string
}

type removeBookEvent struct {
	Title string
}

type library struct {
	books map[string]*book
	ch    chan<- pkg.IDomainEvent
}

func (l *library) AddBook(title string) {
	l.SendDomainEvent(&addBookEvent{Title: title})
}

func (l *library) RemoveBook(title string) {
	l.SendDomainEvent(&removeBookEvent{Title: title})
}

func (l *library) GetBook(title string) *book {
	return l.books[title]
}

func (l *library) SetDomainEventsChannel(ch chan<- pkg.IDomainEvent) {
	l.ch = ch
}

func (l *library) SendDomainEvent(event pkg.IDomainEvent) {
	l.ch <- event
}

func (l *library) handleAddBookEvent(event *addBookEvent) {
	l.books[event.Title] = &book{title: event.Title}
}

func (l *library) handleRemoveBookEvent(event *removeBookEvent) {
	delete(l.books, event.Title)
}

func (l *library) Handle(event pkg.IDomainEvent) {
	switch e := event.(type) {
	case *addBookEvent:
		l.handleAddBookEvent(e)
	case *removeBookEvent:
		l.handleRemoveBookEvent(e)
	}
}

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

type book struct {
	title string
}

func (b *book) validate() error {
	if b.title == "" {
		return errors.New("title is required")
	}
	return nil
}

func newLibrary() *library {
	return &library{
		books: make(map[string]*book),
	}
}

func TestAggregateRoot(t *testing.T) {
	t.Run("should handle valid book operations", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Clean Code")
		time.Sleep(50 * time.Millisecond)

		// Assert
		assert.NotNil(t, lib.GetBook("Clean Code"))
		assert.Equal(t, 0, ar.Version)
		assert.Len(t, ar.GetChanges(), 1)
	})

	t.Run("should detect invalid book state", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		_, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("")
		time.Sleep(50 * time.Millisecond)

		// Assert
		select {
		case err := <-errCh:
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "title is required")
		default:
			t.Fatal("expected error for invalid book")
		}
	})

	t.Run("should maintain version across operations", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		// Act & Assert
		lib.AddBook("Book1")
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 0, ar.Version)

		lib.AddBook("Book2")
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 1, ar.Version)

		lib.RemoveBook("Book1")
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 2, ar.Version)
	})

	t.Run("should load history correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		history := []pkg.IDomainEvent{
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
	})

	t.Run("should handle sequential add and remove operations", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		expectedEvents := 500

		// Act
		for i := 0; i < 250; i++ {
			lib.AddBook(fmt.Sprintf("Book%d", i))
			lib.RemoveBook(fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Equal(c, ar.GetProcessedEventsCount(), expectedEvents)
		}, 2*time.Second, 10*time.Millisecond, "not all events were processed")

		assert.Equal(t, 499, ar.Version)
		assert.Len(t, ar.GetChanges(), expectedEvents)
		for i := 0; i < 250; i++ {
			assert.Nil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential add operations", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		expectedEvents := 500

		// Act
		for i := 0; i < 500; i++ {
			lib.AddBook(fmt.Sprintf("Book%d", i))
		}

		// Assert
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Equal(c, ar.GetProcessedEventsCount(), expectedEvents)
		}, 2*time.Second, 10*time.Millisecond, "not all events were processed")

		assert.Equal(t, 499, ar.Version)
		assert.Len(t, ar.GetChanges(), expectedEvents)
		for i := 0; i < 500; i++ {
			assert.NotNil(t, lib.GetBook(fmt.Sprintf("Book%d", i)))
		}
	})

	t.Run("should handle 500 sequential remove operations", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
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
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Equal(c, ar.GetProcessedEventsCount(), expectedEvents)
		}, 2*time.Second, 10*time.Millisecond, "not all events were processed")

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
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Test Book")
		// time.Sleep(50 * time.Millisecond)
		initialVersion := ar.Version
		ar.ClearChanges()

		// Assert
		assert.NotNil(t, lib.GetBook("Test Book"))
		assert.Empty(t, ar.GetChanges())
		assert.Equal(t, initialVersion, ar.Version)
	})

	t.Run("should handle valid book operations with context", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Clean Code")

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
		errCh := make(chan error, 1)
		pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)

		// Act
		lib.AddBook("Test Book")
		cancel() // Cancel context

		// Assert - try to send after context cancelled
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
		pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
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
}

func BenchmarkAggregateRoot(b *testing.B) {
	b.Run("Sequential AddBook", func(b *testing.B) {
		ctx, _ := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		// defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lib.AddBook("Book")
		}

		for ar.GetProcessedEventsCount() < b.N {
			runtime.Gosched()
		}
	})

	b.Run("Load History", func(b *testing.B) {
		ctx, _ := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		// defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		history := []pkg.IDomainEvent{
			&addBookEvent{Title: "Book1"},
			&removeBookEvent{Title: "Book1"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ar.Load(history)
		}

		// for ar.GetProcessedEventsCount() < b.N {
		// 	runtime.Gosched()
		// }
	})
}

func BenchmarkParallelAggregateRoot(b *testing.B) {
	b.Run("Parallel AddBook", func(b *testing.B) {
		ctx, _ := context.WithTimeout(context.Background(), 10000*time.Millisecond)
		// defer cancel()

		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](ctx, lib, errCh)
		defer closer()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				lib.AddBook("Book")
			}
		})

		for ar.GetProcessedEventsCount() < b.N {
			runtime.Gosched()
		}
	})
}
