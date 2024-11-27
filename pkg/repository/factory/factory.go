//go:generate mockgen -source=factory.go -destination=../../mocks/repository_factory_mock.go -package=mocks
package factory

import (
	"github.com/sosalejandro/ddd-golang/pkg/repository"
)

// RepositoryFactory interface for creating GenericRepository instances
type RepositoryFactory interface {
	CreateRepository() repository.AggregateRepositoryInterface
}
