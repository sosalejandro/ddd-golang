package pkg

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	"github.com/stretchr/testify/assert"
)

// DummyEvent is a mock implementation for testing.
type DummyEvent struct {
	Value string
}

func (de *DummyEvent) SomeMethod() {}

// DummyEvent1 is a mock event for testing.
type DummyEvent1 struct {
	Value string
}

func (de1 *DummyEvent1) SomeMethod() {}

// DummyEvent2 is a mock event for testing.
type DummyEvent2 struct {
	Value string
}

func (de2 *DummyEvent2) SomeMethod() {}

// DummyEvent3 is a mock event for testing.
type DummyEvent3 struct {
	Value string
}

func (de3 *DummyEvent3) SomeMethod() {}

// DummyEvent4 is a mock event for testing.
type DummyEvent4 struct {
	Value string
}

func (de4 *DummyEvent4) SomeMethod() {}

// DummyEvent5 is a mock event for testing.
type DummyEvent5 struct {
	Value string
}

func (de5 *DummyEvent5) SomeMethod() {}

// DummyEvent6 is a mock event for testing.
type DummyEvent6 struct {
	Value string
}

func (de6 *DummyEvent6) SomeMethod() {}

// DummyEvent7 is a mock event for testing.
type DummyEvent7 struct {
	Value string
}

func (de7 *DummyEvent7) SomeMethod() {}

// DummyEvent8 is a mock event for testing.
type DummyEvent8 struct {
	Value string
}

func (de8 *DummyEvent8) SomeMethod() {}

// DummyEvent9 is a mock event for testing.
type DummyEvent9 struct {
	Value string
}

func (de9 *DummyEvent9) SomeMethod() {}

// DummyEvent10 is a mock event for testing.
type DummyEvent10 struct {
	Value string
}

func (de10 *DummyEvent10) SomeMethod() {}

// MultiValueDummyEvent1 is a mock event for testing.
type MultiValueDummyEvent1 struct {
	Value1 string
}

// MultiValueDummyEvent2 is a mock event for testing.
type MultiValueDummyEvent2 struct {
	Value1 string
	Value2 string
}

// MultiValueDummyEvent3 is a mock event for testing.
type MultiValueDummyEvent3 struct {
	Value1 string
	Value2 string
	Value3 string
}

// MultiValueDummyEvent4 is a mock event for testing.
type MultiValueDummyEvent4 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
}

// MultiValueDummyEvent5 is a mock event for testing.
type MultiValueDummyEvent5 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
}

// MultiValueDummyEvent6 is a mock event for testing.
type MultiValueDummyEvent6 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
	Value6 string
}

// MultiValueDummyEvent7 is a mock event for testing.
type MultiValueDummyEvent7 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
	Value6 string
	Value7 string
}

// MultiValueDummyEvent8 is a mock event for testing.
type MultiValueDummyEvent8 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
	Value6 string
	Value7 string
	Value8 string
}

// MultiValueDummyEvent9 is a mock event for testing.
type MultiValueDummyEvent9 struct {
	Value1 string
	Value2 string
	Value3 string
	Value4 string
	Value5 string
	Value6 string
	Value7 string
	Value8 string
	Value9 string
}

// MultiValueDummyEvent10 is a mock event for testing.
type MultiValueDummyEvent10 struct {
	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Value5  string
	Value6  string
	Value7  string
	Value8  string
	Value9  string
	Value10 string
}

// TestNewEventManager tests the creation of a new EventManager.
func TestNewEventManager(t *testing.T) {
	t.Parallel()

	// Arrange
	var em *EventManager

	// Act
	em = NewEventManager()

	// Assert
	assert.NotNil(t, em.registry, "Expected registry to be initialized")
}

// TestRegisterEvents tests registering events in EventManager.
func TestRegisterEvents(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}

	// Act
	em.RegisterEvents(dummyEvent)

	// Assert
	exists := em.Exists("DummyEvent")
	assert.True(t, exists, "Expected DummyEvent to be registered")
}

// TestUnmarshalEvent tests unmarshaling an event.
func TestUnmarshalEvent(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	em.RegisterEvents(dummyEvent)
	Metadata, _ := json.Marshal(dummyEvent)

	// Act
	event, err := em.UnmarshalEvent("DummyEvent", Metadata)

	// Assert
	assert.NoError(t, err)
	unmarshaledEvent, ok := event.(*DummyEvent)
	assert.True(t, ok, "Expected type *DummyEvent")
	assert.Equal(t, "test value", unmarshaledEvent.Value, "Expected event Value to match")
}

// TestUnmarshalEventUnknownType tests unmarshaling an unknown event type.
func TestUnmarshalEventUnknownType(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()
	Metadata := []byte(`{"some_field": "some_value"}`)

	// Act
	event, err := em.UnmarshalEvent("UnknownEvent", Metadata)

	// Assert
	assert.Error(t, err, "Expected error for unknown event type")
	assert.Nil(t, event, "Expected event to be nil for unknown event type")
}

// TestUnmarshalEventInvalidJSON tests unmarshaling with invalid JSON Metadata.
func TestUnmarshalEventInvalidJSON(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	em.RegisterEvents(dummyEvent)
	invalidData := []byte(`{"invalid_json": `) // Malformed JSON

	// Act
	event, err := em.UnmarshalEvent("DummyEvent", invalidData)

	// Assert
	assert.Error(t, err, "Expected error for invalid JSON Metadata")
	assert.Nil(t, event, "Expected event to be nil for invalid JSON Metadata")
}

// TestUnmarshalMultipleEvents validates that multiple events are correctly unmarshalled into their respective types.
func TestUnmarshalMultipleEvents(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()

	// Register multiple events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{Value: "value1"},
		&DummyEvent2{Value: "value2"},
		&DummyEvent3{Value: "value3"},
		&DummyEvent4{Value: "value4"},
		&DummyEvent5{Value: "value5"},
		&DummyEvent6{Value: "value6"},
		&DummyEvent7{Value: "value7"},
		&DummyEvent8{Value: "value8"},
		&DummyEvent9{Value: "value9"},
		&DummyEvent10{Value: "value10"},
	}
	em.RegisterEvents(events...)

	// Act & Assert
	// Marshal and unmarshal each event
	for _, event := range events {
		Metadata, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Failed to marshal %T: %v", event, err)
		}

		eventType := reflect.TypeOf(event).Elem().Name()
		unmarshalled, err := em.UnmarshalEvent(eventType, Metadata)
		if err != nil {
			t.Fatalf("Failed to unmarshal %s: %v", eventType, err)
		}

		assert.Equal(t, reflect.TypeOf(event), reflect.TypeOf(unmarshalled), "Expected type %T, got %T", event, unmarshalled)
	}
}

// TestUnmarshalEventsTableDriven tests unmarshalling multiple events using a table-driven approach.
func TestUnmarshalEventsTableDriven(t *testing.T) {
	t.Parallel()

	// Arrange
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{Value: "value1"},
		&DummyEvent2{Value: "value2"},
		&DummyEvent3{Value: "value3"},
		&DummyEvent4{Value: "value4"},
		&DummyEvent5{Value: "value5"},
		&DummyEvent6{Value: "value6"},
		&DummyEvent7{Value: "value7"},
		&DummyEvent8{Value: "value8"},
		&DummyEvent9{Value: "value9"},
		&DummyEvent10{Value: "value10"},
	}
	em.RegisterEvents(events...)

	testCases := []struct {
		name      string
		numEvents int
	}{
		{"UnmarshalSingle", 1},
		{"UnmarshalPair", 2},
		{"UnmarshalFive", 5},
		{"UnmarshalTen", 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			selectedEvents := events[:tc.numEvents]

			for _, event := range selectedEvents {
				Metadata, err := json.Marshal(event)
				if err != nil {
					t.Fatalf("Failed to marshal %T: %v", event, err)
				}

				// Act
				eventType := reflect.TypeOf(event).Elem().Name()
				unmarshalled, err := em.UnmarshalEvent(eventType, Metadata)
				if err != nil {
					t.Fatalf("Failed to unmarshal %s: %v", eventType, err)
				}

				// Assert
				assert.Equal(t, reflect.TypeOf(event), reflect.TypeOf(unmarshalled), "Expected type %T, got %T", event, unmarshalled)

				// Cleanup
				// em.PutEventBack(eventType, unmarshalled)
				// Removed as pooling is no longer used
			}
		})
	}
}

// BenchmarkNewEventManagerNonParallel benchmarks the creation of EventManager without parallelism.
func BenchmarkNewEventManagerNonParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEventManager()
	}
}

// BenchmarkNewEventManagerParallel benchmarks the creation of EventManager with parallelism.
func BenchmarkNewEventManagerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewEventManager()
		}
	})
}

// BenchmarkRegisterEventsNonParallel benchmarks registering events without parallelism.
func BenchmarkRegisterEventsNonParallel(b *testing.B) {
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	for i := 0; i < b.N; i++ {
		em.RegisterEvents(dummyEvent)
	}
}

// BenchmarkRegisterEventsParallel benchmarks registering events with parallelism.
func BenchmarkRegisterEventsParallel(b *testing.B) {
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			em.RegisterEvents(dummyEvent)
		}
	})
}

// BenchmarkUnmarshalEventNonParallel benchmarks unmarshaling an event without parallelism.
func BenchmarkUnmarshalEventNonParallel(b *testing.B) {
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	em.RegisterEvents(dummyEvent)
	Metadata, _ := json.Marshal(dummyEvent)
	for i := 0; i < b.N; i++ {
		event, err := em.UnmarshalEvent("DummyEvent", Metadata)
		if err != nil {
			b.Fatal(err)
		}
		_ = event // Prevent unused variable
	}
}

// BenchmarkUnmarshalEventParallel benchmarks unmarshaling an event with parallelism.
func BenchmarkUnmarshalEventParallel(b *testing.B) {
	em := NewEventManager()
	var dummyEvent aggregate.DomainEventInterface = &DummyEvent{Value: "test value"}
	em.RegisterEvents(dummyEvent)
	Metadata, _ := json.Marshal(dummyEvent)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			event, err := em.UnmarshalEvent("DummyEvent", Metadata)
			if err != nil {
				b.Fatal(err)
			}
			_ = event // Prevent unused variable
		}
	})
}

// BenchmarkUnmarshalMultipleEventsNonParallel benchmarks unmarshalling multiple event types without parallelism.
func BenchmarkUnmarshalMultipleEventsNonParallel(b *testing.B) {
	em := NewEventManager()

	// Register multiple events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{Value: "value1"},
		&DummyEvent2{Value: "value2"},
		&DummyEvent3{Value: "value3"},
		&DummyEvent4{Value: "value4"},
		&DummyEvent5{Value: "value5"},
		&DummyEvent6{Value: "value6"},
		&DummyEvent7{Value: "value7"},
		&DummyEvent8{Value: "value8"},
		&DummyEvent9{Value: "value9"},
		&DummyEvent10{Value: "value10"},
	}

	em.RegisterEvents(events...)

	var marshalledEvents []EventPayload
	for _, event := range events {
		Metadata, err := json.Marshal(event)
		if err != nil {
			b.Fatalf("Failed to marshal %T: %v", event, err)
		}
		eventType := reflect.TypeOf(event).Elem().Name()
		marshalledEvents = append(marshalledEvents, EventPayload{
			EventType: eventType,
			Data:      Metadata,
		})
	}

	for i := 0; i < b.N; i++ {
		for _, me := range marshalledEvents {
			event, err := em.UnmarshalEvent(me.EventType, me.Data)
			if err != nil {
				b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
			}
			_ = event // Prevent unused variable
		}
	}
}

// BenchmarkUnmarshalMultipleEventsParallel benchmarks unmarshalling multiple event types with parallelism.
func BenchmarkUnmarshalMultipleEventsParallel(b *testing.B) {
	em := NewEventManager()

	// Register multiple events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{Value: "value1"},
		&DummyEvent2{Value: "value2"},
		&DummyEvent3{Value: "value3"},
		&DummyEvent4{Value: "value4"},
		&DummyEvent5{Value: "value5"},
		&DummyEvent6{Value: "value6"},
		&DummyEvent7{Value: "value7"},
		&DummyEvent8{Value: "value8"},
		&DummyEvent9{Value: "value9"},
		&DummyEvent10{Value: "value10"},
	}

	em.RegisterEvents(events...)

	var marshalledEvents []EventPayload
	for _, event := range events {
		Metadata, err := json.Marshal(event)
		if err != nil {
			b.Fatalf("Failed to marshal %T: %v", event, err)
		}
		eventType := reflect.TypeOf(event).Elem().Name()
		marshalledEvents = append(marshalledEvents, EventPayload{
			EventType: eventType,
			Data:      Metadata,
		})
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, me := range marshalledEvents {
				event, err := em.UnmarshalEvent(me.EventType, me.Data)
				if err != nil {
					b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
				}
				_ = event // Prevent unused variable
			}
		}
	})
}

func BenchmarkUnmarshalEventsTableDrivenNonParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{},
		&DummyEvent2{},
		&DummyEvent3{},
		&DummyEvent4{},
		&DummyEvent5{},
		&DummyEvent6{},
		&DummyEvent7{},
		&DummyEvent8{},
		&DummyEvent9{},
		&DummyEvent10{},
	}
	em.RegisterEvents(events...)

	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"UnmarshalSingle", generateEventsArray(events[0], b)},
		{"UnmarshalPair", generateEventsArray(events[1], b)},
		{"UnmarshalFive", generateEventsArray(events[4], b)},
		{"UnmarshalTen", generateEventsArray(events[9], b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name+"NonParallel", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, me := range tc.events {
					event, err := em.UnmarshalEvent(me.EventType, me.Data)
					if err != nil {
						b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
					}
					_ = event // Prevent unused variable
				}
			}
		})
	}
}

func BenchmarkUnmarshalEventsTableDrivenParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{},
		&DummyEvent2{},
		&DummyEvent3{},
		&DummyEvent4{},
		&DummyEvent5{},
		&DummyEvent6{},
		&DummyEvent7{},
		&DummyEvent8{},
		&DummyEvent9{},
		&DummyEvent10{},
	}
	em.RegisterEvents(events...)

	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"UnmarshalSingle", generateEventsArray(events[0], b)},
		{"UnmarshalPair", generateEventsArray(events[1], b)},
		{"UnmarshalFive", generateEventsArray(events[4], b)},
		{"UnmarshalTen", generateEventsArray(events[9], b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name+"Parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for _, me := range tc.events {
						event, err := em.UnmarshalEvent(me.EventType, me.Data)
						if err != nil {
							b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
						}
						_ = event // Prevent unused variable
					}
				}
			})
		})
	}
}

func BenchmarkMultiValueEventTableDrivenNonParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&MultiValueDummyEvent1{},
		&MultiValueDummyEvent2{},
		&MultiValueDummyEvent3{},
		&MultiValueDummyEvent4{},
		&MultiValueDummyEvent5{},
		&MultiValueDummyEvent6{},
		&MultiValueDummyEvent7{},
		&MultiValueDummyEvent8{},
		&MultiValueDummyEvent9{},
		&MultiValueDummyEvent10{},
	}
	em.RegisterEvents(events...)

	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"MultiValueDummyEvent1", generateEventsArray(&MultiValueDummyEvent1{Value1: "value1"}, b)},
		{"MultiValueDummyEvent2", generateEventsArray(&MultiValueDummyEvent2{Value1: "value1", Value2: "value2"}, b)},
		{"MultiValueDummyEvent3", generateEventsArray(&MultiValueDummyEvent3{Value1: "value1", Value2: "value2", Value3: "value3"}, b)},
		{"MultiValueDummyEvent4", generateEventsArray(&MultiValueDummyEvent4{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4"}, b)},
		{"MultiValueDummyEvent5", generateEventsArray(&MultiValueDummyEvent5{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5"}, b)},
		{"MultiValueDummyEvent6", generateEventsArray(&MultiValueDummyEvent6{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6"}, b)},
		{"MultiValueDummyEvent7", generateEventsArray(&MultiValueDummyEvent7{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7"}, b)},
		{"MultiValueDummyEvent8", generateEventsArray(&MultiValueDummyEvent8{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8"}, b)},
		{"MultiValueDummyEvent9", generateEventsArray(&MultiValueDummyEvent9{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8", Value9: "value9"}, b)},
		{"MultiValueDummyEvent10", generateEventsArray(&MultiValueDummyEvent10{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8", Value9: "value9", Value10: "value10"}, b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name, func(b *testing.B) {
			for _, me := range tc.events {
				event, err := em.UnmarshalEvent(me.EventType, me.Data)
				if err != nil {
					b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
				}
				_ = event // Prevent unused variable
			}
		})
	}
}

func BenchmarkMultiValueEventTableDrivenParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&MultiValueDummyEvent1{},
		&MultiValueDummyEvent2{},
		&MultiValueDummyEvent3{},
		&MultiValueDummyEvent4{},
		&MultiValueDummyEvent5{},
		&MultiValueDummyEvent6{},
		&MultiValueDummyEvent7{},
		&MultiValueDummyEvent8{},
		&MultiValueDummyEvent9{},
		&MultiValueDummyEvent10{},
	}
	em.RegisterEvents(events...)

	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"MultiValueDummyEvent1", generateEventsArray(&MultiValueDummyEvent1{Value1: "value1"}, b)},
		{"MultiValueDummyEvent2", generateEventsArray(&MultiValueDummyEvent2{Value1: "value1", Value2: "value2"}, b)},
		{"MultiValueDummyEvent3", generateEventsArray(&MultiValueDummyEvent3{Value1: "value1", Value2: "value2", Value3: "value3"}, b)},
		{"MultiValueDummyEvent4", generateEventsArray(&MultiValueDummyEvent4{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4"}, b)},
		{"MultiValueDummyEvent5", generateEventsArray(&MultiValueDummyEvent5{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5"}, b)},
		{"MultiValueDummyEvent6", generateEventsArray(&MultiValueDummyEvent6{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6"}, b)},
		{"MultiValueDummyEvent7", generateEventsArray(&MultiValueDummyEvent7{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7"}, b)},
		{"MultiValueDummyEvent8", generateEventsArray(&MultiValueDummyEvent8{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8"}, b)},
		{"MultiValueDummyEvent9", generateEventsArray(&MultiValueDummyEvent9{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8", Value9: "value9"}, b)},
		{"MultiValueDummyEvent10", generateEventsArray(&MultiValueDummyEvent10{Value1: "value1", Value2: "value2", Value3: "value3", Value4: "value4", Value5: "value5", Value6: "value6", Value7: "value7", Value8: "value8", Value9: "value9", Value10: "value10"}, b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for _, me := range tc.events {
						event, err := em.UnmarshalEvent(me.EventType, me.Data)
						if err != nil {
							b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
						}
						_ = event // Prevent unused variable
					}
				}
			})
		})
	}
}

// Helper function to generate an array of marshalled events
func generateEventsArray(event aggregate.DomainEventInterface, b *testing.B) []EventPayload {
	var marshalledEvents []EventPayload
	eventType := reflect.TypeOf(event).Elem().Name()
	for i := 0; i < 10; i++ { // Generate 10 identical events
		Metadata, err := json.Marshal(event)
		if err != nil {
			b.Fatalf("Failed to marshal %T: %v", event, err)
		}
		marshalledEvents = append(marshalledEvents, EventPayload{
			EventType: eventType,
			Data:      Metadata,
		})
	}
	return marshalledEvents
}

func BenchmarkUnmarshalMultipleSameEventsNonParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{},
		&DummyEvent2{},
		&DummyEvent3{},
		&DummyEvent4{},
		&DummyEvent5{},
		&DummyEvent6{},
		&DummyEvent7{},
		&DummyEvent8{},
		&DummyEvent9{},
		&DummyEvent10{},
	}
	em.RegisterEvents(events...)

	// Define test cases for each DummyEvent type
	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"DummyEvent1", generateEventsArray(&DummyEvent1{Value: "value1"}, b)},
		{"DummyEvent2", generateEventsArray(&DummyEvent2{Value: "value2"}, b)},
		{"DummyEvent3", generateEventsArray(&DummyEvent3{Value: "value3"}, b)},
		{"DummyEvent4", generateEventsArray(&DummyEvent4{Value: "value4"}, b)},
		{"DummyEvent5", generateEventsArray(&DummyEvent5{Value: "value5"}, b)},
		{"DummyEvent6", generateEventsArray(&DummyEvent6{Value: "value6"}, b)},
		{"DummyEvent7", generateEventsArray(&DummyEvent7{Value: "value7"}, b)},
		{"DummyEvent8", generateEventsArray(&DummyEvent8{Value: "value8"}, b)},
		{"DummyEvent9", generateEventsArray(&DummyEvent9{Value: "value9"}, b)},
		{"DummyEvent10", generateEventsArray(&DummyEvent10{Value: "value10"}, b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name, func(b *testing.B) {
			for _, me := range tc.events {
				event, err := em.UnmarshalEvent(me.EventType, me.Data)
				if err != nil {
					b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
				}
				_ = event // Prevent unused variable
			}
		})
	}
}

func BenchmarkUnmarshalMultipleSameEventsParallel(b *testing.B) {
	em := NewEventManager()

	// Register all events
	events := []aggregate.DomainEventInterface{
		&DummyEvent1{},
		&DummyEvent2{},
		&DummyEvent3{},
		&DummyEvent4{},
		&DummyEvent5{},
		&DummyEvent6{},
		&DummyEvent7{},
		&DummyEvent8{},
		&DummyEvent9{},
		&DummyEvent10{},
	}
	em.RegisterEvents(events...)

	// Define test cases for each DummyEvent type
	testCases := []struct {
		name   string
		events []EventPayload
	}{
		{"DummyEvent1", generateEventsArray(&DummyEvent1{Value: "value1"}, b)},
		{"DummyEvent2", generateEventsArray(&DummyEvent2{Value: "value2"}, b)},
		{"DummyEvent3", generateEventsArray(&DummyEvent3{Value: "value3"}, b)},
		{"DummyEvent4", generateEventsArray(&DummyEvent4{Value: "value4"}, b)},
		{"DummyEvent5", generateEventsArray(&DummyEvent5{Value: "value5"}, b)},
		{"DummyEvent6", generateEventsArray(&DummyEvent6{Value: "value6"}, b)},
		{"DummyEvent7", generateEventsArray(&DummyEvent7{Value: "value7"}, b)},
		{"DummyEvent8", generateEventsArray(&DummyEvent8{Value: "value8"}, b)},
		{"DummyEvent9", generateEventsArray(&DummyEvent9{Value: "value9"}, b)},
		{"DummyEvent10", generateEventsArray(&DummyEvent10{Value: "value10"}, b)},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		b.Run(tc.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for _, me := range tc.events {
						event, err := em.UnmarshalEvent(me.EventType, me.Data)
						if err != nil {
							b.Fatalf("Failed to unmarshal %s: %v", me.EventType, err)
						}
						_ = event // Prevent unused variable
					}
				}
			})
		})
	}
}
