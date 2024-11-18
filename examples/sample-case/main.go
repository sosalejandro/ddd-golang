package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

func main() {
	exampleBaseCaseUsage()
	exampleLoadUsage()
}

func exampleBaseCaseUsage() {
	// Create a library.
	lib := newLibrary()

	// Create an error channel with buffer.
	errCh := make(chan error, 10)

	// Create a context with cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create an aggregate root.
	ar, _ := aggregate.NewAggregateRoot(ctx, lib, errCh)
	// defer closeChannel()

	// Start a goroutine to handle errors.
	go func() {
		for err := range errCh {
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	// Add and remove books.
	lib.AddBook("The Hobbit")
	lib.RemoveBook("The Hobbit")

	// Wait until processed events count is 2 or context is canceled.
	if err := ar.WaitForEvents(ctx, 2); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Example base case usage completed.")
}

func exampleLoadUsage() {
	// Create a library.
	lib := newLibrary()

	// Create an error channel with buffer.
	errCh := make(chan error, 10)

	// Create a context with cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create an aggregate root.
	ar, _ := aggregate.NewAggregateRoot(ctx, lib, errCh)

	// Start a goroutine to handle errors.
	go func() {
		for err := range errCh {
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	history := []aggregate.IDomainEvent{
		&addBookEvent{Title: "The Hobbit"},
		&removeBookEvent{Title: "The Hobbit"},
		&addBookEvent{Title: "Harry Potter and the Philosopher's Stone"},
		&removeBookEvent{Title: "Harry Potter and the Philosopher's Stone"},
		&addBookEvent{Title: "Harry Potter and the Chamber of Secrets"},
		&removeBookEvent{Title: "Harry Potter and the Chamber of Secrets"},
		&addBookEvent{Title: "Harry Potter and the Prisoner of Azkaban"},
		&removeBookEvent{Title: "Harry Potter and the Prisoner of Azkaban"},
	}

	// Load the history.
	ar.Load(history)

	if err := ar.WaitForEvents(ctx, len(history)); err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Processed events count: %d\n", ar.GetProcessedEventsCount())
	fmt.Println("Example load usage completed.")
	for title := range lib.books {
		fmt.Println("-", title)
	}

	// Example of adding another book
	lib.AddBook("1984")
	if err := ar.WaitForEvents(ctx, len(history)+1); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Additional event processed.")
}
