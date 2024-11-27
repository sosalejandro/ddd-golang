package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
)

type EventManager struct {
	registry sync.Map
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) RegisterEvents(events ...aggregate.DomainEventInterface) {
	for _, event := range events {
		eventType := reflect.TypeOf(event).Elem().Name()
		em.registry.Store(eventType, reflect.TypeOf(event).Elem())
	}
}

// Exists checks if an event type exists in the registry.
func (em *EventManager) Exists(eventType string) bool {
	_, ok := em.registry.Load(eventType)
	return ok
}

func (em *EventManager) UnmarshalEvent(eventType string, data []byte) (interface{}, error) {
	eventTypeReflect, ok := em.registry.Load(eventType)
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	event := reflect.New(eventTypeReflect.(reflect.Type)).Interface()
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}
