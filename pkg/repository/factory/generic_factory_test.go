package factory_test

import (
	"testing"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/mocks"
	"github.com/sosalejandro/ddd-golang/pkg/repository"
	"github.com/sosalejandro/ddd-golang/pkg/repository/factory"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type GenericRepositoryFactorySuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockRepoFactory   *mocks.MockRepositoryFactory
	mockAggregate     *mocks.MockAggregateRootInterface
	eventManager      *pkg.EventManager
	genericRepository *repository.GenericRepository[aggregate.AggregateRootInterface]
	genericFactory    *factory.GenericRepositoryFactory[aggregate.AggregateRootInterface]
}

func TestGenericRepositoryFactorySuite(t *testing.T) {
	suite.Run(t, new(GenericRepositoryFactorySuite))
}

func (s *GenericRepositoryFactorySuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepoFactory = mocks.NewMockRepositoryFactory(s.ctrl)
	s.mockAggregate = mocks.NewMockAggregateRootInterface(s.ctrl)
	s.eventManager = pkg.NewEventManager()
	s.genericFactory = factory.NewGenericRepositoryFactory(s.mockRepoFactory, func() aggregate.AggregateRootInterface { return s.mockAggregate }, s.eventManager)
}

func (s *GenericRepositoryFactorySuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GenericRepositoryFactorySuite) TestCreateGenericRepository() {
	s.mockRepoFactory.EXPECT().CreateRepository().Times(1)

	repo := s.genericFactory.CreateGenericRepository()

	s.NotNil(repo)
}
