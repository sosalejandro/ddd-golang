package pkg

import "github.com/google/uuid"

type EventPayload struct {
	AggregateID uuid.UUID
	Timestamp   int64
	EventType   string
	EventID     int // Use aggregate version as EventID
	Data        []byte
}
