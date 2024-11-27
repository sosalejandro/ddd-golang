package memory

import "github.com/sosalejandro/ddd-golang/pkg/repository"

type MemoryRepositoryFactory struct{}

func NewMemoryRepositoryFactory() *MemoryRepositoryFactory {
	return &MemoryRepositoryFactory{}
}

func (f *MemoryRepositoryFactory) CreateRepository() repository.AggregateRepositoryInterface {
	return NewMemoryRepository()
}
