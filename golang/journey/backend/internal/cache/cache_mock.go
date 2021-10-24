package cache

import (
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetUniqueJourneys() []Journey {
	args := m.Called()
	return args.Get(0).([]Journey)
}

func (m *MockCache) StartJourney(characterId string, startId, destinationId uint16) {
	m.Called(characterId, startId, destinationId)
}

func (m *MockCache) Movement(characterId string, x, y uint16) error {
	args := m.Called(characterId, x, y)
	return args.Error(0)
}

func (m *MockCache) ReachedDestination(characterId string, destinationId uint16) error {
	args := m.Called(characterId, destinationId)
	return args.Error(0)
}
