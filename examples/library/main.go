package main

import (
	"context"
	"fmt"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

func main() {
	exampleBaseCaseUsage()
	exampleLoadUsage()
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

	// Load the history.
	if err := lib.Load(ctx, history); err != nil {
		fmt.Println("Error loading history:", err)
	}

	fmt.Println("Example load usage completed.")
	for title := range lib.books {
		fmt.Println("-", title)
	}

	// Example of adding another book
	if err := lib.AddBook(ctx, "1984"); err != nil {
		fmt.Println("Error adding book:", err)
	}
	fmt.Println("Additional event processed.")
}
