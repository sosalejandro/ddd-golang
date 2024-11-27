# Event Manager in Go

This package provides an Event Manager implementation for handling domain events within a Go application, adhering to Domain-Driven Design (DDD) principles. It facilitates the registration, existence check, and unmarshalling of events, ensuring a scalable and maintainable event-driven architecture.

## Overview

The `EventManager` is responsible for managing domain events within the application. It allows registering event types, checking for their existence, and unmarshalling event data into corresponding event structures. This ensures that events are consistently handled and processed across different parts of the system.

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.18 or higher installed on your machine.
- [Git](https://git-scm.com/) for cloning the repository.

### Installation

1. **Clone the Repository**
```bash
git clone https://github.com/sosalejandro/ddd-golang.git
cd ddd-golang/pkg/event-manager
```

2. **Initialize Go Modules**
```bash
go mod tidy
```

## Usage

### Initializing the Event Manager

Create a new instance of the `EventManager`:

```go
import "github.com/sosalejandro/ddd-golang/pkg/event-manager"

em := event_manager.NewEventManager()
```

### Registering Events

Register your domain events with the Event Manager:

```go
type BookAddedEvent struct {
    Title string
}

type BookRemovedEvent struct {
    Title string
}

em.RegisterEvents(&BookAddedEvent{}, &BookRemovedEvent{})
```

### Checking Event Existence

Verify if an event type is registered:

```go
exists := em.Exists("BookAddedEvent")
if exists {
    fmt.Println("Event type exists.")
} else {
    fmt.Println("Event type does not exist.")
}
```

### Unmarshalling Events

Convert JSON data into the corresponding event structure:

```go
data := []byte(`{"Title": "1984"}`)
event, err := em.UnmarshalEvent("BookAddedEvent", data)
if err != nil {
    fmt.Println("Error unmarshalling event:", err)
} else {
    fmt.Printf("Unmarshalled Event: %+v\n", event)
}
```

### Using EventManager as a Singleton

The `EventManager` can be instantiated as a singleton to manage all events across the service. This approach ensures that event registrations are centralized, reducing redundancy and facilitating easier maintenance.

```go
var singletonEventManager = event_manager.NewEventManager()

func GetEventManager() *event_manager.EventManager {
    return singletonEventManager
}
```

Register all events during the application's initialization phase:

```go
func init() {
    em := GetEventManager()
    em.RegisterEvents(&BookAddedEvent{}, &BookRemovedEvent{})
}
```

This singleton instance is optimized for scenarios with numerous read operations and minimal write operations, primarily during initialization.

### Using EventManager per Aggregate

Alternatively, you can instantiate separate `EventManager` instances for each aggregate. This approach encapsulates event management within specific aggregates, promoting modularity and separation of concerns.

```go
func NewLibraryEventManager() *event_manager.EventManager {
    em := event_manager.NewEventManager()
    em.RegisterEvents(&BookAddedEvent{}, &BookRemovedEvent{})
    return em
}
```

Each aggregate interacts with its own `EventManager`, allowing for tailored event handling and reduced coupling between different parts of the system.

### Optimization for Reads Over Writes

The `EventManager` is designed to be highly efficient for read operations, making it ideal for applications where event retrieval is frequent and event registrations are infrequent. It employs concurrent-safe structures like `sync.Map` to facilitate fast, lock-free reads. 

- **Reads:** Optimized using concurrent data structures to allow multiple goroutines to access the registry without contention.
- **Writes:** Primarily occur during the initialization phase when registering events. Since writes are minimal and rarely occur after setup, the performance impact is negligible.

This optimization strategy ensures that the `EventManager` can handle high read loads efficiently while maintaining thread safety.

## Example

```go
package main

import (
    "fmt"
    "github.com/sosalejandro/ddd-golang/pkg/event-manager"
)

func main() {
    em := event_manager.NewEventManager()

    // Register events
    em.RegisterEvents(&BookAddedEvent{}, &BookRemovedEvent{})

    // Check event existence
    if em.Exists("BookAddedEvent") {
        fmt.Println("BookAddedEvent is registered.")
    }

    // Unmarshal event data
    data := []byte(`{"Title": "Brave New World"}`)
    event, err := em.UnmarshalEvent("BookAddedEvent", data)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Received Event: %+v\n", event)
    }
}

type BookAddedEvent struct {
    Title string
}

type BookRemovedEvent struct {
    Title string
}
```

## Handling Unknown Events

If an event type is not registered, the `UnmarshalEvent` method will return an error. Ensure that all necessary event types are registered before attempting to unmarshal.

```go
event, err := em.UnmarshalEvent("UnknownEvent", data)
if err != nil {
    fmt.Println("Error:", err) // Outputs: unknown event type: UnknownEvent
}
```

## Dependencies

- [Google UUID](https://github.com/google/uuid) for generating unique identifiers.

## Testing

Unit tests are provided to ensure the correctness of the Event Manager operations.

```bash
go test ./pkg/event-manager
```

## Conclusion

The `EventManager` package offers a robust solution for managing domain events in Go applications. By registering, checking, and unmarshalling events efficiently, it supports the implementation of a clean and maintainable event-driven architecture. Whether used as a singleton or per aggregate, its optimization for read-heavy operations makes it suitable for a wide range of applications.