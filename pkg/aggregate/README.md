# Aggregate Root Example in Go

This repository demonstrates the implementation of the `AggregateRoot` pattern in Go, adhering to Domain-Driven Design (DDD) principles. The example encapsulates a simple library system where books can be added and removed through domain events.

## Overview

The `AggregateRoot` is central to managing domain events and maintaining the state of the aggregate (in this case, the `Library`). It ensures that all business rules are enforced and the state remains consistent throughout operations.

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.18 or higher installed on your machine.
- [Git](https://git-scm.com/) for cloning the repository.

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/sosalejandro/ddd-golang.git
   cd ddd-golang/pkg/aggregate
   ```

2. **Initialize Go Modules**

   ```bash
   go mod tidy
   ```

3. **Run the Example**

   Navigate to the example directory and execute the application:

   ```bash
   cd examples/library
   go run main.go
   ```

## Usage

### Instantiating AggregateRoot

To utilize the `AggregateRoot`, follow these steps:

1. **Define Your Aggregate Root Entity**

   The `Library` struct serves as the aggregate root, managing a collection of books.

2. **Create Domain Events**

   Define domain events such as `addBookEvent` and `removeBookEvent` to represent significant state changes.

3. **Initialize AggregateRoot**

   Instantiate the `AggregateRoot` and associate it with the `Library` entity.

   ```go
   lib := newLibrary()
   ctx := context.Background()

   if err := lib.AddBook(ctx, "1984"); err != nil {
       fmt.Println("Error adding book:", err)
   }
   ```

### Performing Operations

#### Adding a Book

```go
if err := lib.AddBook(ctx, "Brave New World"); err != nil {
    fmt.Println("Error adding book:", err)
}
```

#### Removing a Book

```go
if err := lib.RemoveBook(ctx, "1984"); err != nil {
    fmt.Println("Error removing book:", err)
}
```

### Loading Event History

Reconstruct the state of the `Library` by loading a history of domain events.

```go
history := []aggregate.RecordedEvent{
    {EventID: 0, Timestamp: time.Now().Unix(), Event: &addBookEvent{Title: "1984"}},
    {EventID: 1, Timestamp: time.Now().Unix(), Event: &removeBookEvent{Title: "1984"}},
}

if err := lib.Load(ctx, history...); err != nil {
    fmt.Println("Error loading history:", err)
}
```

### Handling Errors and Validation

The `AggregateRoot` enforces business rules. Attempting to add a book with an empty title will result in an error.

```go
err := lib.AddBook(ctx, "")
if err != nil {
    fmt.Println("Error:", err) // Outputs: title is required
}
```

## Aggregate Root and Entity Internals

The `AggregateRoot` is responsible for handling domain events and maintaining the state of the `Library`. It uses channels for asynchronous event processing.

```go
package aggregate

type AggregateRoot struct {
    changes []RecordedEvent
    Version int
}

// ApplyDomainEvent applies a domain event and updates the aggregate.
func (ar *AggregateRoot) ApplyDomainEvent(ctx context.Context, event DomainEventInterface, handle func(ctx context.Context, e DomainEventInterface) error) error {
    // Event handling logic
}
```

The `Library` entity implements methods to add and remove books, handling domain events through the `AggregateRoot`.

```go
package main

type library struct {
    Books map[string]*book
    *aggregate.AggregateRoot
}

func (l *library) AddBook(ctx context.Context, title string) error {
    return l.ApplyDomainEvent(ctx, &addBookEvent{Title: title}, l.Handle)
}

func (l *library) Handle(ctx context.Context, event aggregate.DomainEventInterface) error {
    switch e := event.(type) {
    case *addBookEvent:
        return l.handleAddBookEvent(e)
    case *removeBookEvent:
        return l.handleRemoveBookEvent(e)
    }
    return nil
}
```

### Internal Communication

The `AggregateRoot` uses channels to process domain events asynchronously, ensuring that state mutations adhere to business rules and maintaining eventual consistency.

## Dependencies

- [Go](https://golang.org/) 1.18+
- [Testify](https://github.com/stretchr/testify) for assertions in tests.

## Testing

Unit tests are provided to ensure the correctness of the aggregate root operations.

```bash
go test ./pkg/aggregate
```

## Benchmarks

Benchmark tests are included to evaluate the performance of the aggregate root.

```bash
go test -bench=.
```

## Conclusion

This example showcases the implementation of the `AggregateRoot` pattern in Go, facilitating clean architecture and separation of concerns. By leveraging Go's interfaces and concurrency features, you can build robust and maintainable domain-driven applications.
