# GenericRepository

The `GenericRepository` is a versatile and reusable component designed to manage aggregates within a Domain-Driven Design (DDD) architecture in Go. It leverages generics to accommodate any aggregate that implements the `AggregateRootInterface`, providing a consistent and scalable approach to data persistence and retrieval.

## Table of Contents

- [GenericRepository](#genericrepository)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Key Components](#key-components)
    - [`AggregateRootInterface`](#aggregaterootinterface)
    - [`AggregateRepositoryInterface`](#aggregaterepositoryinterface)
    - [`EventManager`](#eventmanager)
  - [Usage](#usage)
    - [Initialization](#initialization)
    - [Saving Events](#saving-events)
    - [Rehydrating Aggregates](#rehydrating-aggregates)
    - [Snapshot Management](#snapshot-management)
  - [Abstractions and Composition](#abstractions-and-composition)
  - [Storage Implementations](#storage-implementations)
  - [Generics and Flexibility](#generics-and-flexibility)
  - [Example](#example)
  - [Conclusion](#conclusion)

## Overview

The `GenericRepository` serves as a foundational component for handling the persistence of aggregate roots. By abstracting the storage mechanism, it enables developers to focus on the domain logic without worrying about the underlying data operations. This repository pattern facilitates a clean separation of concerns, promoting maintainability and scalability in complex applications.

## Key Components

### `AggregateRootInterface`

The `AggregateRootInterface` defines the contract that any aggregate root must adhere to. It ensures that the repository can interact with aggregates uniformly.

```go
package aggregate

// AggregateRootInterface is an interface that represents an aggregate root.
type AggregateRootInterface interface {
	Id() uuid.UUID
	SetId(uuid.UUID)
	GetChanges() []RecordedEvent
	Load(ctx context.Context, events ...RecordedEvent) error
	Deserialize([]byte) error
	Serialize() ([]byte, error)
}
```

### `AggregateRepositoryInterface`

The `AggregateRepositoryInterface` abstracts the data storage operations required by the `GenericRepository`. Any storage implementation (e.g., SQL, NoSQL, in-memory) must implement this interface to be compatible.

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    event_manager "github.com/sosalejandro/ddd-golang/pkg/event-manager"
)

type AggregateRepositoryInterface interface {
    SaveEvents(context.Context, []event_manager.EventPayload) error
    GetAggregateEvents(context.Context, uuid.UUID) ([]event_manager.EventPayload, error)
    LoadSnapshot(context.Context, uuid.UUID) (*Snapshot, error)
    SaveSnapshot(context.Context, *Snapshot) error
    Close() error
}
```

### `EventManager`

The `EventManager` handles the registration and unmarshalling of domain events. It ensures that events are correctly processed and stored, maintaining the integrity of the event-driven architecture.

## Usage

### Initialization

To create a new instance of `GenericRepository`, provide it with an `EventManager`, an implementation of `AggregateRepositoryInterface`, and a factory function that returns a new instance of the aggregate.

```go
import (
    event_manager "github.com/sosalejandro/ddd-golang/pkg/event-manager"
    "github.com/sosalejandro/ddd-golang/pkg/repository"
    "github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

em := event_manager.NewEventManager()
repo := NewConcreteRepository() // Your implementation of AggregateRepositoryInterface
factory := func() aggregate.AggregateRootInterface {
    return &YourAggregate{}
}

genericRepo := repository.NewGenericRepository(em, repo, factory)
```

### Saving Events

The `SaveEvents` method persists all uncommitted changes (events) of an aggregate.

```go
ctx := context.Background()
aggregate := genericRepo.Load(ctx, aggregateID)
aggregate.SomeMethod()
err := genericRepo.SaveEvents(ctx, aggregate)
if err != nil {
    // Handle error
}
```

### Rehydrating Aggregates

To reconstruct an aggregate's state, use the `Rehydrate` method, which loads events from the repository and applies them to the aggregate.

```go
ctx := context.Background()
aggregate, err := genericRepo.Rehydrate(ctx, aggregateID)
if err != nil {
    // Handle error
}
```

### Snapshot Management

Snapshots capture the state of an aggregate at a specific point in time, improving performance by reducing the number of events that need to be replayed.

- **Loading a Snapshot:**

  ```go
  aggregate, err := genericRepo.Load(ctx, aggregateID)
  if err != nil {
      // Handle error
  }
  ```

- **Saving a Snapshot:**

  ```go
  err := genericRepo.SaveSnapshot(ctx, aggregate)
  if err != nil {
      // Handle error
  }
  ```

## Abstractions and Composition

The `GenericRepository` leverages abstraction and composition to decouple the domain logic from the persistence mechanism. By depending on interfaces (`AggregateRepositoryInterface` and `AggregateRootInterface`), it promotes flexibility and testability. This design allows developers to switch out or modify storage implementations without altering the core business logic.

## Storage Implementations

Any storage backend (e.g., PostgreSQL, MongoDB, in-memory) can be integrated with the `GenericRepository` by implementing the `AggregateRepositoryInterface`. This ensures that the repository can persist and retrieve aggregates seamlessly, regardless of the underlying storage technology.

```go
type InMemoryRepository struct {
    // Implementation details
}

func (r *InMemoryRepository) SaveEvents(ctx context.Context, events []event_manager.EventPayload) error {
    // Save events to in-memory store
}

func (r *InMemoryRepository) GetAggregateEvents(ctx context.Context, id uuid.UUID) ([]event_manager.EventPayload, error) {
    // Retrieve events from in-memory store
}

func (r *InMemoryRepository) LoadSnapshot(ctx context.Context, id uuid.UUID) (*Snapshot, error) {
    // Load snapshot from in-memory store
}

func (r *InMemoryRepository) SaveSnapshot(ctx context.Context, snapshot *Snapshot) error {
    // Save snapshot to in-memory store
}

func (r *InMemoryRepository) Close() error {
    // Close any resources if necessary
}
```

## Generics and Flexibility

By utilizing Go's generics, the `GenericRepository` can accommodate any aggregate type that implements the `AggregateRootInterface`. This reduces boilerplate code and enhances reusability across different aggregates within the application.

```go
type GenericRepository[T aggregate.AggregateRootInterface] struct {
    em      *event_manager.EventManager
    repo    AggregateRepositoryInterface
    factory func() T
}
```

## Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/sosalejandro/ddd-golang/pkg/aggregate"
    event_manager "github.com/sosalejandro/ddd-golang/pkg/event-manager"
    "github.com/sosalejandro/ddd-golang/pkg/repository"
)

// YourAggregate implements AggregateRootInterface
type YourAggregate struct {
    *aggregate.AggregateRoot
    // Your fields
}

func main() {
    em := event_manager.NewEventManager()
    repo := NewInMemoryRepository() // Your implementation
    factory := func() aggregate.AggregateRootInterface {
        return &YourAggregate{}
    }

    genericRepo := repository.NewGenericRepository(em, repo, factory)

    ctx := context.Background()
    aggregateID := uuid.New()

    // Rehydrate aggregate
    aggregate, err := genericRepo.Rehydrate(ctx, aggregateID)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Perform operations on aggregate
    // ...

    // Save events
    err = genericRepo.SaveEvents(ctx, aggregate)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Save snapshot
    err = genericRepo.SaveSnapshot(ctx, aggregate)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Close repository
    err = genericRepo.Close()
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

## Conclusion

The `GenericRepository` provides a robust and flexible foundation for managing aggregates in a DDD-based Go application. By abstracting storage concerns and leveraging generics, it facilitates the development of scalable and maintainable systems. Whether you're integrating with SQL databases, NoSQL stores, or in-memory solutions, the `GenericRepository` adapts seamlessly, empowering developers to focus on delivering rich domain functionalities.