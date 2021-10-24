package metrics

import (
	"github.com/stretchr/testify/mock"
)

type MockMetrics struct {
	mock.Mock
}

func (m *MockMetrics) LogRequest() {
	m.Called()
}

func (m *MockMetrics) LogJourney() {
	m.Called()
}

func (m *MockMetrics) Close() {
	m.Called()
}
