package factory

import (
	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

// ComposedRepositoryFactory is a factory that composes a RepositoryFactory with an EventManager
type ComposedRepositoryFactory[T aggregate.AggregateRootInterface] struct {
	repoFactory      RepositoryFactory
	aggregateFactory func() T
	eventManager     *pkg.EventManager
}

func NewComposedRepositoryFactory[T aggregate.AggregateRootInterface](
	repoFactory RepositoryFactory,
	aggregateFactory func() T,
	eventManager *pkg.EventManager,
) *ComposedRepositoryFactory[T] {
	return &ComposedRepositoryFactory[T]{
		repoFactory:      repoFactory,
		aggregateFactory: aggregateFactory,
		eventManager:     eventManager,
	}
}

func (f *ComposedRepositoryFactory[T]) CreateGenericRepository() *repository.GenericRepository[T] {
	aggregateFuncs := f.repoFactory.CreateRepository()

	return repository.NewGenericRepository[T](f.eventManager, aggregateFuncs, f.aggregateFactory)
}
