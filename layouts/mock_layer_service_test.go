package layouts

import (
	"github.com/stretchr/testify/mock"
)

type MockLayerService struct {
	mock.Mock
}

func (m *MockLayerService) Put(id, name string, layer []byte) error {
	args := m.Called(id, name, layer)
	return args.Error(0)
}

func (m *MockLayerService) Get(id, name string) ([]byte, string, error) {
	args := m.Called(id, name)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func (m *MockLayerService) Delete(id, name string) error {
	args := m.Called(id, name)
	return args.Error(0)
}
