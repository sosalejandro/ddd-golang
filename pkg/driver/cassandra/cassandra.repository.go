package cassandra

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

const AggregateEventsTable = "aggregate_events"
const AggregateEventCountTable = "aggregate_event_count"
const AggregateEventsVersionTable = "aggregate_events_version"

// AuditRepository defines the read methods for auditing aggregates.
type AuditRepository interface {
	GetEventCount(ctx context.Context, aggregateID uuid.UUID) (int, error)
	GetLatestEventVersion(ctx context.Context, aggregateID uuid.UUID) (int, error)
	ValidateAggregateConsistency(ctx context.Context, aggregateID uuid.UUID) (bool, error)
	GetAllAggregateIDs(ctx context.Context) ([]uuid.UUID, error)
}

// CassandraConfiguration holds the configuration for Cassandra.
type CassandraConfiguration struct {
	ClusterHosts      []string
	Keyspace          string
	EventsTable       string
	SnapshotsTable    string
	LatestEventTable  string
	EventCountTable   string
	EventVersionTable string
}

// AggregateFuncsCassandra implements the AggregateFuncsInterface for Cassandra.
type AggregateFuncsCassandra struct {
	session           *gocql.Session
	keyspace          string
	eventsTable       string
	snapshotsTable    string
	latestEventTable  string
	eventCountTable   string
	eventVersionTable string
	clusterHosts      []string
}

// Reconnect attempts to re-establish the Cassandra session.
func (a *AggregateFuncsCassandra) Reconnect() error {
	cluster := gocql.NewCluster(a.clusterHosts...)
	cluster.Keyspace = a.keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to reconnect to Cassandra: %w", err)
	}

	a.session = session
	return nil
}

// NewAggregateFuncsCassandra initializes a new Cassandra-based AggregateFuncs.
func NewAggregateFuncsCassandra(clusterHosts []string, keyspace, latestEventTable, eventsTable, snapshotsTable, eventCountTable, eventVersionTable string) (*AggregateFuncsCassandra, error) {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	return &AggregateFuncsCassandra{
		session:           session,
		keyspace:          keyspace,
		eventsTable:       eventsTable,
		snapshotsTable:    snapshotsTable,
		latestEventTable:  latestEventTable,
		eventCountTable:   eventCountTable,
		eventVersionTable: eventVersionTable,
		clusterHosts:      clusterHosts,
	}, nil
}

// SaveEvents implements the SaveEvents method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) SaveEvents(ctx context.Context, events []pkg.EventPayload) error {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return err
		}
	}
	if len(events) == 0 {
		return nil
	}
	// iterateEvent Iterates over the events within a batch operation inserting into the events table and the aggregate_events_version table
	var iterateEvent func(ctx context.Context, batch *gocql.Batch, event pkg.EventPayload)
	// saveLatestEventQuery Updates the aggregate_latest_event table with the latest event
	var saveLatestEventQuery func(ctx context.Context, batch *gocql.Batch, latestEvent pkg.EventPayload)
	// saveEventCountQuery Updates the aggregate_event_count table
	var saveEventCountQuery func(ctx context.Context, batch *gocql.Batch, eventCount int, aggregateID gocql.UUID)

	iterateEvent = func(ctx context.Context, batch *gocql.Batch, event pkg.EventPayload) {
		query := fmt.Sprintf(`INSERT INTO %s (aggregate_id, event_id, event_type, data, timestamp) VALUES (?, ?, ?, ?, ?)`, a.eventsTable)
		batch.Query(query, gocql.UUID(event.AggregateID), gocql.UUID(event.EventID), event.EventType, event.Data, event.Timestamp)

		// Insert into aggregate_events_version table
		versionQuery := fmt.Sprintf(`INSERT INTO %s (aggregate_id, event_id, event_version, timestamp) VALUES (?, ?, ?, ?)`, a.eventVersionTable)
		batch.Query(versionQuery, gocql.UUID(event.AggregateID), gocql.UUID(event.EventID), event.Version, event.Timestamp)
	}

	// saveLatestEventQuery Updates the aggregate_latest_event table with the latest event
	saveLatestEventQuery = func(ctx context.Context, batch *gocql.Batch, latestEvent pkg.EventPayload) {
		latestEventQuery := fmt.Sprintf(`UPDATE %s SET latest_event_id = ?, timestamp = ? WHERE aggregate_id = ?`, a.latestEventTable)
		batch.Query(latestEventQuery, gocql.UUID(latestEvent.EventID), latestEvent.Timestamp, gocql.UUID(latestEvent.AggregateID))
	}

	// saveEventCountQuery Updates the aggregate_event_count table
	saveEventCountQuery = func(ctx context.Context, batch *gocql.Batch, eventCount int, aggregateID gocql.UUID) {
		eventCountQuery := fmt.Sprintf(`UPDATE %s SET event_count = ? WHERE aggregate_id = ?`, a.eventCountTable)
		batch.Query(eventCountQuery, eventCount, aggregateID)
	}

	batch := a.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for _, event := range events {
		iterateEvent(ctx, batch, event)
	}

	// Update the aggregate_latest_event table with the latest event
	latestEvent := events[len(events)-1]
	saveLatestEventQuery(ctx, batch, latestEvent)

	saveEventCountQuery(ctx, batch, len(events), gocql.UUID(latestEvent.AggregateID))

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := a.session.ExecuteBatch(batch); err != nil {
			return fmt.Errorf("failed to execute batch: %w", err)
		}
	}

	return nil
}

// GetAggregateEvents implements the GetAggregateEvents method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) GetAggregateEvents(ctx context.Context, id uuid.UUID) ([]pkg.EventPayload, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return nil, err
		}
	}
	var eventsPayload []pkg.EventPayload
	var eventID gocql.UUID
	var eventType string
	var timestamp int64
	var data []byte
	var version int

	query := fmt.Sprintf(`SELECT event_id, event_type, data, timestamp FROM %s WHERE aggregate_id = ? ORDER BY timestamp ASC`, a.eventsTable)

	iter := a.session.Query(query, gocql.UUID(id)).WithContext(ctx).Iter()
	for iter.Scan(&eventID, &eventType, &data, &timestamp, &version) {
		eventsPayload = append(eventsPayload, pkg.EventPayload{
			AggregateID: id,
			EventID:     uuid.UUID(eventID),
			EventType:   eventType,
			Data:        data,
			Timestamp:   timestamp,
			Version:     version,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return eventsPayload, nil
}

// LoadSnapshot implements the LoadSnapshot method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID) (*repository.Snapshot, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return nil, err
		}
	}
	snapshot := new(repository.Snapshot)

	query := fmt.Sprintf(`SELECT aggregate_id, version, data, timestamp FROM %s WHERE aggregate_id = ? ORDER BY timestamp DESC LIMIT 1`, a.snapshotsTable)
	iter := a.session.Query(query, gocql.UUID(aggregateID)).WithContext(ctx).Iter()

	// Check if the snapshot exists
	if !iter.Scan(&snapshot.AggregateID, &snapshot.Version, &snapshot.Data, &snapshot.Timestamp, &snapshot.EventID) {
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("snapshot not found for aggregate ID: %s", aggregateID)
	}

	// Convert gocql.UUID to uuid.UUID
	snapshot.AggregateID = uuid.UUID(snapshot.AggregateID)
	snapshot.EventID = uuid.UUID(snapshot.EventID)

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// SaveSnapshot implements the SaveSnapshot method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) SaveSnapshot(ctx context.Context, snapshot *repository.Snapshot) error {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return err
		}
	}
	batch := a.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)

	query := fmt.Sprintf(`INSERT INTO %s (aggregate_id, data, version, timestamp) VALUES (?, ?, ?, ?)`, a.snapshotsTable)
	batch.Query(query, gocql.UUID(snapshot.AggregateID), snapshot.Data, snapshot.Version, snapshot.Timestamp)

	// Update the aggregate_latest_event table
	latestEventQuery := fmt.Sprintf(`INSERT INTO %s (aggregate_id, latest_event_id) VALUES (?, ?)`, a.latestEventTable)
	batch.Query(latestEventQuery, gocql.UUID(snapshot.AggregateID), gocql.UUID(snapshot.EventID))

	if err := a.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	return nil
}

// Close terminates the Cassandra session.
func (a *AggregateFuncsCassandra) Close() error {
	if a.session != nil {
		a.session.Close()
	}
	// Additional cleanup if needed
	return nil
}

// GetEventCount retrieves the event count for an aggregate.
func (a *AggregateFuncsCassandra) GetEventCount(ctx context.Context, aggregateID uuid.UUID) (int, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return 0, err
		}
	}
	var eventCount int
	query := fmt.Sprintf(`SELECT event_count FROM %s WHERE aggregate_id = ?`, a.eventCountTable)
	if err := a.session.Query(query, gocql.UUID(aggregateID)).WithContext(ctx).Scan(&eventCount); err != nil {
		return 0, fmt.Errorf("failed to retrieve event count: %w", err)
	}
	return eventCount, nil
}

// GetLatestEventVersion retrieves the latest event version for an aggregate.
func (a *AggregateFuncsCassandra) GetLatestEventVersion(ctx context.Context, aggregateID uuid.UUID) (int, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return 0, err
		}
	}
	var latestVersion int
	query := fmt.Sprintf(`SELECT event_version FROM %s WHERE aggregate_id = ? ORDER BY event_version DESC LIMIT 1`, a.eventVersionTable)
	if err := a.session.Query(query, gocql.UUID(aggregateID)).WithContext(ctx).Scan(&latestVersion); err != nil {
		return 0, fmt.Errorf("failed to retrieve latest event version: %w", err)
	}
	return latestVersion, nil
}

// ValidateAggregateConsistency checks if the aggregate is consistent.
func (a *AggregateFuncsCassandra) ValidateAggregateConsistency(ctx context.Context, aggregateID uuid.UUID) (bool, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return false, err
		}
	}
	eventCount, err := a.GetEventCount(ctx, aggregateID)
	if err != nil {
		return false, err
	}

	latestVersion, err := a.GetLatestEventVersion(ctx, aggregateID)
	if err != nil {
		return false, err
	}

	if eventCount != latestVersion {
		return false, fmt.Errorf("inconsistent state: event count (%d) does not match latest version (%d)", eventCount, latestVersion)
	}

	return true, nil
}

// GetAllAggregateIDs retrieves all aggregate IDs.
func (a *AggregateFuncsCassandra) GetAllAggregateIDs(ctx context.Context) ([]uuid.UUID, error) {
	if a.session == nil {
		if err := a.Reconnect(); err != nil {
			return nil, err
		}
	}
	var aggregateIDs []gocql.UUID
	query := fmt.Sprintf(`SELECT aggregate_id FROM %s`, a.eventCountTable)
	iter := a.session.Query(query).WithContext(ctx).Iter()
	for iter.Scan(&aggregateIDs) {
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	// Convert gocql.UUID to uuid.UUID
	var result []uuid.UUID
	for _, id := range aggregateIDs {
		result = append(result, uuid.UUID(id))
	}

	return result, nil
}

// func (s *Scheduler) ValidateAggregates(ctx context.Context) {
//     aggregates, err := s.repo.GetAllAggregateIDs(ctx)
//     if err != nil {
//         log.Fatalf("Failed to retrieve aggregate IDs: %v", err)
//     }

//     for _, aggregateID := range aggregates {
//         isValid, err := s.repo.ValidateAggregateConsistency(ctx, aggregateID)
//         if err != nil {
//             log.Printf("Failed to validate aggregate %s: %v", aggregateID, err)
//             continue
//         }

//         if (!isValid) {
//             log.Printf("Inconsistent aggregate detected: %s", aggregateID)
//             // Handle inconsistency (e.g., alert, fix, etc.)
//         } else {
//             log.Printf("Aggregate %s is consistent", aggregateID)
//         }
//     }
// }
