package cassandra

import (
	"fmt"

	gocql "github.com/gocql/gocql"
)

// SeedDatabase creates the necessary keyspace, tables, and schema in Cassandra.
func SeedDatabase(clusterHosts []string, keyspace string) error {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to Cassandra: %w", err)
	}
	defer session.Close()

	// Create keyspace
	keyspaceQuery := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
		'class': 'SimpleStrategy',
		'replication_factor': '1'
	};`, keyspace)
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Switch to the new keyspace
	cluster.Keyspace = keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to keyspace: %w", err)
	}
	defer session.Close()

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS aggregate_events (
			aggregateId UUID,
			eventId UUID,
			eventType TEXT,
			data BLOB,
			timestamp TIMESTAMP,
			PRIMARY KEY (aggregateId, timestamp)
		) WITH CLUSTERING ORDER BY (timestamp ASC);`,
		`CREATE TABLE IF NOT EXISTS aggregate_latest_event (
			aggregate_id UUID,
			latest_event_id UUID,
			timestamp TIMESTAMP,
			PRIMARY KEY (aggregate_id)
		);`,
		`CREATE TABLE IF NOT EXISTS aggregate_event_count (
			aggregate_id UUID,
			event_count INT,
			PRIMARY KEY (aggregate_id)
		);`,
		`CREATE TABLE IF NOT EXISTS aggregate_events_version (
			aggregate_id UUID,
			event_id UUID,
			event_version INT,
			timestamp TIMESTAMP,
			PRIMARY KEY (aggregate_id, event_version)
		) WITH CLUSTERING ORDER BY (event_version ASC);`,
		`CREATE TABLE IF NOT EXISTS snapshots (
			aggregate_id UUID,
			data BLOB,
			version INT,
			timestamp TIMESTAMP,
			PRIMARY KEY (aggregate_id, timestamp)
		) WITH CLUSTERING ORDER BY (timestamp DESC);`,
	}

	for _, query := range queries {
		if err := session.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
