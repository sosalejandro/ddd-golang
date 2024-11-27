package factory_test

import (
	"testing"

	"github.com/sosalejandro/ddd-golang/pkg/aggregate"
	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/mocks"
	"github.com/sosalejandro/ddd-golang/pkg/repository/factory"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RepositoryFactoryBuilderSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockRepoFactory *mocks.MockRepositoryFactory
	mockAggregate   *mocks.MockAggregateRootInterface
	eventManager    *pkg.EventManager
	builder         *factory.RepositoryFactoryBuilder
}

func TestRepositoryFactoryBuilderSuite(t *testing.T) {
	suite.Run(t, new(RepositoryFactoryBuilderSuite))
}

func (s *RepositoryFactoryBuilderSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepoFactory = mocks.NewMockRepositoryFactory(s.ctrl)
	s.mockAggregate = mocks.NewMockAggregateRootInterface(s.ctrl)
	aggFactory := func() aggregate.AggregateRootInterface { return s.mockAggregate }
	s.eventManager = pkg.NewEventManager()
	s.builder = factory.NewRepositoryFactoryBuilder().
		WithConfig("").
		WithEventManager(s.eventManager).
		WithAggregateFactory(aggFactory)
}

func (s *RepositoryFactoryBuilderSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RepositoryFactoryBuilderSuite) TestBuild_Success() {
	// ...existing code...
	s.mockRepoFactory.EXPECT().CreateRepository().Return(s.mockRepoFactory.CreateRepository()).Times(1)

	// ...existing code...
	repoFactory, err := s.builder.Build()

	// ...existing code...
	s.NoError(err)
	s.NotNil(repoFactory)
}

func (s *RepositoryFactoryBuilderSuite) TestBuild_MissingConfig() {
	// ...existing code...
	s.builder = factory.NewRepositoryFactoryBuilder().
		WithEventManager(s.eventManager)

	// ...existing code...
	repoFactory, err := s.builder.Build()

	// ...existing code...
	s.Error(err)
	s.Nil(repoFactory)
}

func (s *RepositoryFactoryBuilderSuite) TestBuild_UnsupportedConfig() {
	// ...existing code...
	s.builder = factory.NewRepositoryFactoryBuilder().
		WithConfig(make(chan string)).
		WithEventManager(s.eventManager).
		WithAggregateFactory(func() aggregate.AggregateRootInterface { return s.mockAggregate })

	// ...existing code...
	repoFactory, err := s.builder.Build()

	// ...existing code...
	s.Error(err)
	s.Nil(repoFactory)
}
