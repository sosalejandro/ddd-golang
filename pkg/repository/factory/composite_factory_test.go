package factory_test

import (
	"testing"

	pkg "github.com/sosalejandro/ddd-golang/pkg/event-manager"
	"github.com/sosalejandro/ddd-golang/pkg/mocks"
	"github.com/sosalejandro/ddd-golang/pkg/repository/factory"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CompositeFactorySuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockRepoFactory  *mocks.MockRepositoryFactory
	mockAggregate    *mocks.MockAggregateRootInterface
	eventManager     *pkg.EventManager
	compositeFactory *factory.ComposedRepositoryFactory[*mocks.MockAggregateRootInterface]
}

func TestCompositeFactorySuite(t *testing.T) {
	suite.Run(t, new(CompositeFactorySuite))
}

func (s *CompositeFactorySuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepoFactory = mocks.NewMockRepositoryFactory(s.ctrl)
	s.mockAggregate = mocks.NewMockAggregateRootInterface(s.ctrl)
	s.eventManager = pkg.NewEventManager()
	aggFactory := func() *mocks.MockAggregateRootInterface { return s.mockAggregate }
	s.compositeFactory = factory.NewComposedRepositoryFactory(s.mockRepoFactory, aggFactory, s.eventManager)
}

func (s *CompositeFactorySuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CompositeFactorySuite) TestCreateCompositeRepository_Success() {
	// ...existing code...
	s.mockRepoFactory.EXPECT().CreateRepository().Times(1)

	// ...existing code...
	compositeRepo := s.compositeFactory.CreateGenericRepository()

	// ...existing code...
	s.NotNil(compositeRepo)
}
