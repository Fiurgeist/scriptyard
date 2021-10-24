package queue

import (
	"github.com/stretchr/testify/mock"
)

type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) Push(msg interface{}) {
	m.Called(msg)
}

func (m *MockQueue) GetChannel() chan interface{} {
	args := m.Called()
	return args.Get(0).(chan interface{})
}

func (m *MockQueue) Close() {
	m.Called()
}
