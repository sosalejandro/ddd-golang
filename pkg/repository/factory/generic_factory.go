package factory

import (
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

// GenericRepositoryFactory is a factory for creating GenericRepository instances
type GenericRepositoryFactory[T aggregate.AggregateRootInterface] struct {
	factory          RepositoryFactory
	aggregateFactory func() T
	eventManager     *pkg.EventManager
}

func NewGenericRepositoryFactory[T aggregate.AggregateRootInterface](factory RepositoryFactory, aggregateFactory func() T, eventManager *pkg.EventManager) *GenericRepositoryFactory[T] {
	return &GenericRepositoryFactory[T]{
		factory:          factory,
		aggregateFactory: aggregateFactory,
		eventManager:     eventManager,
	}
}

func (grf *GenericRepositoryFactory[T]) CreateGenericRepository() *repository.GenericRepository[T] {
	aggregateFuncs := grf.factory.CreateRepository()

	return repository.NewGenericRepository[T](grf.eventManager, aggregateFuncs, grf.aggregateFactory)
}
