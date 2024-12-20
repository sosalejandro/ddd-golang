package pkg

import "github.com/google/uuid"

type EventPayload struct {
	AggregateID uuid.UUID
	Timestamp   int64
	EventType   string
	EventID     uuid.UUID // Change to UUID
	Data        []byte
	Version     int
}
