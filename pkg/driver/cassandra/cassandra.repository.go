package cassandra

import (
	"context"
	"fmt"

	gocql "github.com/gocql/gocql"
	"github.com/google/uuid"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

const AggregateEventsTable = "aggregate_events"

// CassandraConfiguration holds the configuration for Cassandra.
type CassandraConfiguration struct {
	ClusterHosts   []string
	Keyspace       string
	AggregateTable string
	EventsTable    string
	SnapshotsTable string
}

// AggregateFuncsCassandra implements the AggregateFuncsInterface for Cassandra.
type AggregateFuncsCassandra struct {
	session        *gocql.Session
	keyspace       string
	aggregateTable string
	eventsTable    string
	snapshotsTable string
}

// NewAggregateFuncsCassandra initializes a new Cassandra-based AggregateFuncs.
func NewAggregateFuncsCassandra(clusterHosts []string, keyspace, aggregateTable, eventsTable, snapshotsTable string) (*AggregateFuncsCassandra, error) {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	return &AggregateFuncsCassandra{
		session:        session,
		keyspace:       keyspace,
		aggregateTable: aggregateTable,
		eventsTable:    eventsTable,
		snapshotsTable: snapshotsTable,
	}, nil
}

// Save implements the Save method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) SaveEvents(ctx context.Context, events []pkg.EventPayload) error {
	batch := a.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for _, event := range events {
		query := fmt.Sprintf(`INSERT INTO %s (id, eventId, eventType, data, timestamp) VALUES (?, ?, ?)`, a.eventsTable)
		batch.Query(query, event.AggregateID, event.EventID, event.EventType, event.Data, event.Timestamp)
	}

	if err := a.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	return nil
}

// GetEvents implements the GetEvents method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) GetAggregateEvents(ctx context.Context, id uuid.UUID) ([]pkg.EventPayload, error) {
	var eventsPayload []pkg.EventPayload
	var eventID int
	var eventType string
	var timestamp int64
	var data []byte

	query := fmt.Sprintf(`SELECT eventId, eventType, data, timestamp FROM %s WHERE id = ? LIMIT 1`, a.eventsTable)

	iter := a.session.Query(query, id).WithContext(ctx).Iter()
	for iter.Scan(&eventID, &eventType, &data, &timestamp) {
		eventsPayload = append(eventsPayload, pkg.EventPayload{
			AggregateID: id,
			EventID:     eventID,
			EventType:   eventType,
			Data:        data,
			Timestamp:   timestamp,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return eventsPayload, nil
}

// LoadSnapshot implements the LoadSnapshot method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID) (*repository.Snapshot, error) {
	snapshot := new(repository.Snapshot)

	query := fmt.Sprintf(`SELECT aggregateId, version, data, timestamp FROM %s WHERE aggregateId = ? LIMIT 1`, a.snapshotsTable)
	iter := a.session.Query(query, aggregateID).WithContext(ctx).Iter()

	if iter.Scan(&snapshot.AggregateID, &snapshot.Version, &snapshot.Data, &snapshot.Timestamp) {
		if err := iter.Close(); err != nil {
			return nil, err
		}

		return snapshot, nil
	}

	return nil, nil
}

// SaveSnapshot implements the SaveSnapshot method of AggregateFuncsInterface.
func (a *AggregateFuncsCassandra) SaveSnapshot(ctx context.Context, snapshot *repository.Snapshot) error {
	query := fmt.Sprintf(`INSERT INTO %s (aggregate_id, data, version, timestamp) VALUES (?, ?, ?)`, a.snapshotsTable)

	if err := a.session.Query(query, snapshot.AggregateID, snapshot.Data, snapshot.Version, snapshot.Timestamp).WithContext(ctx).Exec(); err != nil {
		return err
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
