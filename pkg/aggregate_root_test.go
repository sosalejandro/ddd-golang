package pkg_test

import (
	"errors"
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
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
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
		lib := newLibrary()
		errCh := make(chan error, 1)
		_, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
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
		lib := newLibrary()
		errCh := make(chan error, 2)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
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
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
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

	// t.Run("should handle concurrent operations", func(t *testing.T) {
	// 	// Arrange
	// 	lib := newLibrary()
	// 	errCh := make(chan error, 10)
	// 	_, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
	// 	defer closer()

	// 	// Act
	// 	go lib.AddBook("Book1")
	// 	go lib.AddBook("Book2")
	// 	go lib.RemoveBook("Book1")

	// 	// Assert
	// 	time.Sleep(200 * time.Millisecond)
	// 	assert.Nil(t, lib.GetBook("Book1"))
	// 	assert.NotNil(t, lib.GetBook("Book2"))
	// })

	t.Run("should clear changes without affecting state", func(t *testing.T) {
		// Arrange
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
		defer closer()

		// Act
		lib.AddBook("Test Book")
		time.Sleep(50 * time.Millisecond)
		initialVersion := ar.Version
		ar.ClearChanges()

		// Assert
		assert.NotNil(t, lib.GetBook("Test Book"))
		assert.Empty(t, ar.GetChanges())
		assert.Equal(t, initialVersion, ar.Version)
	})
}

func BenchmarkAggregateRoot(b *testing.B) {
	b.Run("Sequential AddBook", func(b *testing.B) {
		lib := newLibrary()
		errCh := make(chan error, b.N)
		_, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
		defer closer()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lib.AddBook("Book")
			time.Sleep(1 * time.Millisecond)
		}
	})

	b.Run("Load History", func(b *testing.B) {
		lib := newLibrary()
		errCh := make(chan error, 1)
		ar, closer := pkg.NewAggregateRoot[pkg.IAggregateRoot](lib, errCh)
		defer closer()

		history := []pkg.IDomainEvent{
			&addBookEvent{Title: "Book1"},
			&removeBookEvent{Title: "Book1"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ar.Load(history)
		}
	})
}
