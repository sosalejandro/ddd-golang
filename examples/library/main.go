package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	"github.com/sosalejandro/ddd-golang/pkg/driver/memory"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

func main() {
	exampleBaseCaseUsage()
	exampleLoadUsage()
	exampleRepositorySaveAndLoad()
}

func exampleBaseCaseUsage() {
	ctx := context.Background()
	// Create a library.
	lib := newLibrary()

	// Add and remove books.
	if err := lib.AddBook(ctx, "The Hobbit"); err != nil {
		fmt.Println("Error adding book:", err)
	}
	if err := lib.RemoveBook(ctx, "The Hobbit"); err != nil {
		fmt.Println("Error removing book:", err)
	}

	fmt.Println("Example base case usage completed.")
}

func exampleLoadUsage() {
	// Create a context.
	ctx := context.Background()

	// Create a library.
	lib := newLibrary()

	history := []aggregate.DomainEventInterface{
		&addBookEvent{Title: "The Hobbit"},
		&removeBookEvent{Title: "The Hobbit"},
		&addBookEvent{Title: "Harry Potter and the Philosopher's Stone"},
		&removeBookEvent{Title: "Harry Potter and the Philosopher's Stone"},
		&addBookEvent{Title: "Harry Potter and the Chamber of Secrets"},
		&removeBookEvent{Title: "Harry Potter and the Chamber of Secrets"},
		&addBookEvent{Title: "Harry Potter and the Prisoner of Azkaban"},
		&removeBookEvent{Title: "Harry Potter and the Prisoner of Azkaban"},
	}

	eventRecords := make([]aggregate.RecordedEvent, len(history))
	for i, event := range history {
		eventRecords = append(eventRecords, aggregate.RecordedEvent{
			Event:     event,
			Version:   i,
			EventID:   uuid.New(),
			Timestamp: time.Now().Unix(),
		})
	}

	// Load the history.
	if err := lib.Load(ctx, eventRecords...); err != nil {
		fmt.Println("Error loading history:", err)
	}

	fmt.Println("Example load usage completed.")
	for title := range lib.Books {
		fmt.Println("-", title)
	}

	// Example of adding another book
	if err := lib.AddBook(ctx, "1984"); err != nil {
		fmt.Println("Error adding book:", err)
	}
	fmt.Println("Additional event processed.")
}

func exampleRepositorySaveAndLoad() {
	// Initialize EventManager
	em := pkg.NewEventManager()

	em.RegisterEvents(
		&addBookEvent{},
		&removeBookEvent{},
	)

	// Create MemoryRepository using factory
	factory := memory.NewMemoryRepositoryFactory()
	memoryRepo := factory.CreateRepository()

	// Create GenericRepository for library aggregate
	repo := repository.NewGenericRepository(em, memoryRepo, newLibrary)

	// Create a new library
	lib := newLibrary()
	lib.SetId(uuid.New())

	// Perform operations
	ctx := context.Background()
	lib.AddBook(ctx, "1984")
	lib.AddBook(ctx, "Brave New World")

	// Save events
	if err := repo.SaveEvents(ctx, lib); err != nil {
		fmt.Println("Error saving events:", err)
		return
	}

	// Save snapshot
	if err := repo.SaveSnapshot(ctx, lib); err != nil {
		fmt.Println("Error saving snapshot:", err)
		return
	}

	fmt.Println("Library saved successfully.")

	// Load library
	lib, err := repo.Load(ctx, lib.Id())
	if err != nil {
		fmt.Println("Error loading library:", err)
		return
	}

	// Display loaded books
	fmt.Printf("Library Changes: %v\n", lib.Books)
	for title := range lib.Books {
		fmt.Println("-", title)
	}

	lib2, err := repo.Rehydrate(ctx, lib.Id())
	if err != nil {
		fmt.Println("Error rehydrating library:", err)
		return
	}

	// Display loaded books
	fmt.Printf("Library 2 Changes: %v\n", lib2.Books)
	for title := range lib2.Books {
		fmt.Println("-", title)
	}

	fmt.Println("Example repository save and load completed.")
}

// func exampleRepositoryLoad() {
// 	// Initialize EventManager
// 	em := pkg.NewEventManager()

// 	// Create MemoryRepository using factory
// 	factory := memory.NewMemoryRepositoryFactory()
// 	memoryRepo := factory.CreateRepository()

// 	// Create GenericRepository for library aggregate
// 	repo := repository.NewGenericRepository[*library](em, memoryRepo)

// 	// Assume we have an aggregate ID
// 	aggregateID := uuid.New()

// 	// Load aggregate
// 	lib, err := repo.Rehydrate(context.Background(), aggregateID)
// 	if err != nil {
// 		fmt.Println("Error loading library:", err)
// 		return
// 	}

// 	// Display loaded books
// 	for title := range lib.books {
// 		fmt.Println("-", title)
// 	}
// }
