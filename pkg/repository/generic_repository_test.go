package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/mocks"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type GenericRepositorySuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockRepository    *mocks.MockAggregateRepositoryInterface
	mockAggregate     *mocks.MockAggregateRootInterface
	eventManager      *pkg.EventManager
	genericRepository *repository.GenericRepository[*mocks.MockAggregateRootInterface]
}

func TestGenericRepositorySuite(t *testing.T) {
	suite.Run(t, new(GenericRepositorySuite))
}

func (s *GenericRepositorySuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepository = mocks.NewMockAggregateRepositoryInterface(s.ctrl)
	s.mockAggregate = mocks.NewMockAggregateRootInterface(s.ctrl)
	s.eventManager = pkg.NewEventManager()
	s.eventManager.RegisterEvents(
		&sampleDomainEvent{},
		&invalidEvent{},
	)

	factory := func() *mocks.MockAggregateRootInterface {
		return s.mockAggregate
	}
	s.genericRepository = repository.NewGenericRepository(s.eventManager, s.mockRepository, factory)
}

func (s *GenericRepositorySuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GenericRepositorySuite) TestSaveEvents() {
	// Setup mock expectations for GetChanges
	expectedChanges := []aggregate.RecordedEvent{
		{
			EventID:   1,
			Timestamp: time.Now().Unix(),
			Event:     &sampleDomainEvent{},
		},
	}

	s.mockAggregate.EXPECT().
		GetChanges().
		Return(expectedChanges).
		Times(1)

	// Setup mock expectations for Id with specific arguments
	aggregateId := uuid.New()
	s.mockAggregate.EXPECT().Id().Return(aggregateId).Times(1)

	// Transform expectedChanges to expected EventPayloads
	expectedEventPayloads := expectedChangesToSave(expectedChanges, aggregateId)

	// Setup mock expectations for SaveEvents with specific arguments
	ctx := context.Background()
	s.mockRepository.EXPECT().
		SaveEvents(ctx, expectedEventPayloads).
		Return(nil).
		Times(1)

	// Call SaveEvents
	err := s.genericRepository.SaveEvents(ctx, s.mockAggregate)

	// Assert no error
	s.NoError(err)
}

func expectedChangesToSave(changes []aggregate.RecordedEvent, aggregateId uuid.UUID) []pkg.EventPayload {
	var payloads []pkg.EventPayload
	for _, change := range changes {
		data, _ := json.Marshal(change.Event)
		payloads = append(payloads, pkg.EventPayload{
			AggregateID: aggregateId,
			Timestamp:   change.Timestamp,
			EventType:   reflect.TypeOf(change.Event).Elem().Name(),
			EventID:     change.EventID,
			Data:        data,
		})
	}
	return payloads
}

func (s *GenericRepositorySuite) TestRehydrate() {
	aggregateID := uuid.New()
	data, _ := json.Marshal(&sampleDomainEvent{})
	expectedEvents := []pkg.EventPayload{
		{
			AggregateID: aggregateID,
			Timestamp:   time.Now().Unix(),
			EventType:   "sampleDomainEvent",
			EventID:     1,
			Data:        data,
		},
	}

	// Setup mock expectations for Id with specific arguments
	s.mockAggregate.EXPECT().SetId(aggregateID).Times(1)

	// Setup mock expectations for GetAggregateEvents with specific arguments
	ctx := context.Background()
	s.mockAggregate.EXPECT().Load(ctx, gomock.Any()).Return(nil).Times(1)
	s.mockRepository.EXPECT().
		GetAggregateEvents(ctx, aggregateID).
		Return(expectedEvents, nil).
		Times(1)

	// Setup mock expectations for Deserialize if applicable
	// Assuming Rehydrate might invoke Deserialize through Load or other methods
	// If not, this can be omitted

	// Call Rehydrate
	aggregate, err := s.genericRepository.Rehydrate(ctx, aggregateID)

	// Assert no error and returned aggregate is correctly rehydrated
	s.NoError(err)
	// Further assertions can be added based on the behavior of Rehydrate
	s.NotNil(aggregate)
}

func (s *GenericRepositorySuite) TestLoad() {
	aggregateID := uuid.New()
	snapshot := &repository.Snapshot{ // Changed to pointer
		AggregateID: aggregateID,
		Data:        []byte("{}"),
		Version:     1,
		Timestamp:   time.Now().Unix(),
	}

	// Setup mock expectations for Id with specific arguments
	s.mockAggregate.EXPECT().SetId(aggregateID).Times(1)
	s.mockAggregate.EXPECT().Id().Return(aggregateID).Times(1)

	// Setup mock expectations for LoadSnapshot with specific arguments
	ctx := context.Background()
	s.mockRepository.EXPECT().
		LoadSnapshot(ctx, aggregateID).
		Return(snapshot, nil).
		Times(1)

	// Setup expectation for Deserialize with specific argument
	s.mockAggregate.EXPECT().
		Deserialize(snapshot.Data).
		Return(nil).
		Times(1)

	// Call Load
	aggregate, err := s.genericRepository.Load(ctx, aggregateID)

	// Assert no error and returned aggregate ID is set
	s.NoError(err)
	s.Equal(aggregateID, aggregate.Id())
}

func (s *GenericRepositorySuite) TestSaveSnapshot() {
	// Setup mock expectations for Serialize with specific arguments
	serializedData := []byte("{}")
	s.mockAggregate.EXPECT().
		Serialize().
		Return(serializedData, nil).
		Times(1)

	// Setup mock expectations for GetChanges with specific arguments (called twice)
	expectedChanges := []aggregate.RecordedEvent{
		{
			EventID:   1,
			Timestamp: time.Now().Unix(),
			Event:     &sampleDomainEvent{},
		},
	}
	s.mockAggregate.EXPECT().
		GetChanges().
		Return(expectedChanges).
		Times(2) // Called twice in SaveSnapshot

	// Setup mock expectations for Id with specific arguments
	expectedID := uuid.New()
	s.mockAggregate.EXPECT().
		Id().
		Return(expectedID).
		Times(1)

	// Transform expectedChanges to expected Snapshot
	expectedSnapshot := &repository.Snapshot{
		AggregateID: expectedID,
		Data:        serializedData,
		Version:     expectedChanges[len(expectedChanges)-1].EventID,
		Timestamp:   time.Now().Unix(),
	}

	// Setup mock expectations for SaveSnapshot with specific arguments
	ctx := context.Background()
	s.mockRepository.EXPECT().
		SaveSnapshot(ctx, expectedSnapshot).
		Return(nil).
		Times(1)

	// Call SaveSnapshot
	err := s.genericRepository.SaveSnapshot(ctx, s.mockAggregate)

	// Assert no error
	s.NoError(err)
}

func (s *GenericRepositorySuite) TestClose() {
	// Setup mock expectations for Close with specific arguments
	s.mockRepository.EXPECT().
		Close().
		Return(nil).
		Times(1)

	// Call Close
	err := s.genericRepository.Close()

	// Assert no error
	s.NoError(err)
}

func (s *GenericRepositorySuite) TestSaveEvents_JSONMarshalError() {
	aggregateID := uuid.New()
	// Setup mock to return an error when json.Marshal is called
	s.mockAggregate.EXPECT().
		GetChanges().
		Return([]aggregate.RecordedEvent{
			{
				EventID:   1,
				Timestamp: time.Now().Unix(),
				Event:     make(chan int), // json.Marshal will fail on chan type
			},
		}).
		Times(1)

	// Mock the repository to expect SaveEvents not to be called due to marshal error

	// Setup mock expectations for Id with specific arguments
	s.mockAggregate.EXPECT().Id().Return(aggregateID).Times(1)

	// Call SaveEvents and assert the error
	ctx := context.Background()
	err := s.genericRepository.SaveEvents(ctx, s.mockAggregate)
	s.Error(err)
	// s.Contains(err.Error(), "invalid event")
}

func (s *GenericRepositorySuite) TestRehydrate_GetAggregateEventsError() {
	aggregateID := uuid.New()
	ctx := context.Background()
	// Setup mock to return an error when GetAggregateEvents is called
	s.mockRepository.EXPECT().
		GetAggregateEvents(ctx, aggregateID).
		Return(nil, errors.New("failed to retrieve events")).
		Times(1)

	// Call Rehydrate and assert the error

	_, err := s.genericRepository.Rehydrate(ctx, aggregateID)
	s.Error(err)
	s.Equal("failed to retrieve events", err.Error())
}

func (s *GenericRepositorySuite) TestLoadSnapshot_DeserializeError() {
	aggregateID := uuid.New()
	snapshot := &repository.Snapshot{
		AggregateID: aggregateID,
		Data:        []byte("invalid json"),
		Version:     1,
		Timestamp:   time.Now().Unix(),
	}
	ctx := context.Background()

	// Setup mock expectations for Id with specific arguments
	s.mockAggregate.EXPECT().SetId(aggregateID).Times(1)

	// Setup mock expectations for LoadSnapshot to return invalid JSON
	s.mockRepository.EXPECT().
		LoadSnapshot(ctx, aggregateID).
		Return(snapshot, nil).
		Times(1)

	// Setup mock expectations for Deserialize to return an error
	s.mockAggregate.EXPECT().
		Deserialize(snapshot.Data).
		Return(errors.New("deserialize error")).
		Times(1)

	// Call Load and assert the error
	_, err := s.genericRepository.Load(ctx, aggregateID)
	s.Error(err)
	s.Equal("deserialize error", err.Error())
}

func (s *GenericRepositorySuite) TestSaveSnapshot_SerializeError() {
	// Setup mock to return an error when Serialize is called
	s.mockAggregate.EXPECT().
		Serialize().
		Return(nil, errors.New("serialize error")).
		Times(1)

	// Call SaveSnapshot and assert the error
	ctx := context.Background()
	err := s.genericRepository.SaveSnapshot(ctx, s.mockAggregate)
	s.Error(err)
	s.Equal("serialize error", err.Error())
}

func (s *GenericRepositorySuite) TestRehydrate_EventManagerUnmarshalError() {
	aggregateID := uuid.New()
	invalidEventPayload := []pkg.EventPayload{
		{
			AggregateID: aggregateID,
			Timestamp:   time.Now().Unix(),
			EventType:   "UnknownEvent",
			EventID:     1,
			Data:        []byte("{}"),
		},
	}

	ctx := context.Background()
	// Setup mock to return events with unknown event type
	s.mockRepository.EXPECT().
		GetAggregateEvents(ctx, aggregateID).
		Return(invalidEventPayload, nil).
		Times(1)

	// Call Rehydrate and assert the error

	_, err := s.genericRepository.Rehydrate(ctx, aggregateID)
	s.Error(err)
	// s.Contains(err.Error(), "unmarshal error")
}

// sampleDomainEvent is a mock implementation of DomainEventInterface for testing purposes.
type sampleDomainEvent struct{}

func (e *sampleDomainEvent) SomeMethod() {}

// invalidEvent is a mock event that causes json.Marshal to fail
type invalidEvent struct{}

func (e *invalidEvent) SomeMethod() {}
