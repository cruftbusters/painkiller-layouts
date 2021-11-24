package v1

import (
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/stretchr/testify/mock"
)

type MockAwaitingLayerService struct {
	mock.Mock
}

func (m *MockAwaitingLayerService) Enqueue(got types.Layout) error {
	args := m.Called(got)
	return args.Error(0)
}

func (m *MockAwaitingLayerService) Dequeue() types.Layout {
	args := m.Called()
	return args.Get(0).(types.Layout)
}
