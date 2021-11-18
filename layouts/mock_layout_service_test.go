package layouts

import (
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/stretchr/testify/mock"
)

type MockLayoutService struct {
	mock.Mock
}

func (m *MockLayoutService) Create(layout Layout) Layout {
	args := m.Called(layout)
	return args.Get(0).(Layout)
}

func (m *MockLayoutService) Get(id string) (Layout, error) {
	args := m.Called(id)
	return args.Get(0).(Layout), args.Error(1)
}

func (m *MockLayoutService) GetAll() []Layout {
	args := m.Called()
	return args.Get(0).([]Layout)
}

func (m *MockLayoutService) GetAllWithNoHeightmap() []Layout {
	args := m.Called()
	return args.Get(0).([]Layout)
}

func (m *MockLayoutService) GetAllWithHeightmapWithoutHillshade() []Layout {
	args := m.Called()
	return args.Get(0).([]Layout)
}

func (m *MockLayoutService) Patch(id string, patch Layout) (Layout, error) {
	args := m.Called(id, patch)
	return args.Get(0).(Layout), args.Error(1)
}

func (m *MockLayoutService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
