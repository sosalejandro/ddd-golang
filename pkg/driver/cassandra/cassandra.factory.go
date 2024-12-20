package cassandra

import "github.com/sosalejandro/ddd-golang/pkg/repository"

// CassandraRepositoryFactory implements RepositoryFactory for Cassandra
type CassandraRepositoryFactory struct {
	config *CassandraConfiguration
}

func NewCassandraRepositoryFactory(config *CassandraConfiguration) *CassandraRepositoryFactory {
	return &CassandraRepositoryFactory{config: config}
}

func (f *CassandraRepositoryFactory) CreateRepository() repository.AggregateRepositoryInterface {
	cassandra, err := NewAggregateFuncsCassandra(
		f.config.ClusterHosts,
		f.config.Keyspace,
		f.config.EventsTable,
		f.config.SnapshotsTable,
		f.config.LatestEventTable,
		f.config.EventCountTable,
		f.config.EventVersionTable,
	)
	if err != nil {
		panic(err) // Handle this better in production
	}

	return cassandra
}
