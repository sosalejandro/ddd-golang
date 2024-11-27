package factory

import (
	"fmt"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	"github.com/sosalejandro/ddd-golang/pkg/driver/cassandra"
	"github.com/sosalejandro/ddd-golang/pkg/driver/memory"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
)

// RepositoryFactoryBuilder is a builder for creating a RepositoryFactory
type RepositoryFactoryBuilder struct {
	config           interface{}
	aggregateFactory func() aggregate.AggregateRootInterface
	eventManager     *pkg.EventManager
}

func NewRepositoryFactoryBuilder() *RepositoryFactoryBuilder {
	return &RepositoryFactoryBuilder{}
}

func (b *RepositoryFactoryBuilder) WithConfig(config interface{}) *RepositoryFactoryBuilder {
	b.config = config
	return b
}

func (b *RepositoryFactoryBuilder) WithAggregateFactory(aggregateFactory func() aggregate.AggregateRootInterface) *RepositoryFactoryBuilder {
	b.aggregateFactory = aggregateFactory
	return b
}

func (b *RepositoryFactoryBuilder) WithEventManager(eventManager *pkg.EventManager) *RepositoryFactoryBuilder {
	b.eventManager = eventManager
	return b
}

func (b *RepositoryFactoryBuilder) Build() (*GenericRepositoryFactory[aggregate.AggregateRootInterface], error) {
	if b.config == nil || b.eventManager == nil {
		return nil, fmt.Errorf("config and event manager must be set before building")
	}

	// Determine which RepositoryFactory to create based on the config type
	var repoFactory RepositoryFactory
	switch cfg := b.config.(type) {
	case *cassandra.CassandraConfiguration:
		repoFactory = cassandra.NewCassandraRepositoryFactory(cfg)
	// Add cases for other configurations (e.g., MongoDBConfiguration)
	case string:
		repoFactory = memory.NewMemoryRepositoryFactory()
	default:
		return nil, fmt.Errorf("unsupported config type: %T", b.config)
	}

	// Return the fully built GenericRepositoryFactory
	return NewGenericRepositoryFactory(repoFactory, b.aggregateFactory, b.eventManager), nil
}
