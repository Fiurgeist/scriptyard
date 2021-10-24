package cache

import (
	"fiurgeist/journey/internal/metrics"
	"fiurgeist/journey/internal/queue"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUniqueJourneys(t *testing.T) {
	cache := NewCache(nil, nil)

	// test empty cache
	require.Equal(t, []Journey{}, cache.GetUniqueJourneys())

	// test filled cache
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}
	cache.journeys["42->23"] = &journey{
		startId:       42,
		destinationId: 23,
		points:        []Point{{X: 1, Y: 2}, {X: 2, Y: 2}},
		isFullyMapped: true,
	}

	gotJourneys := cache.GetUniqueJourneys()
	require.Equal(t, 2, len(gotJourneys))
	require.Contains(t, gotJourneys, Journey{Id: "23->42", Points: nil})
	require.Contains(
		t,
		gotJourneys,
		Journey{Id: "42->23", Points: []Point{{X: 1, Y: 2}, {X: 2, Y: 2}}},
	)
}

func TestStartJourney(t *testing.T) {
	mockQueue := &queue.MockQueue{}
	cache := NewCache(nil, mockQueue)

	expectedMsg := queue.NewJourney{StartId: 23, DestinationId: 42}
	mockQueue.On("Push", expectedMsg).Return()

	// first characterJourney for the journey
	cache.StartJourney("character1", 23, 42)

	// add characterJourney and journey
	expectedCharacterJourney1 := &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedJourney1 := &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}
	require.Equal(t, 1, len(cache.characterJourneys))
	require.Equal(t, expectedCharacterJourney1, cache.characterJourneys["character1"])
	require.Equal(t, 1, len(cache.journeys))
	require.Equal(t, expectedJourney1, cache.journeys["23->42"])

	// another characterJourney of a different character but for the same route points
	cache.StartJourney("character2", 23, 42)

	// only add characterJourney
	expectedCharacterJourney2 := &characterJourney{characterId: "character2", startId: 23, destinationId: 42}
	require.Equal(t, 2, len(cache.characterJourneys))
	require.Equal(t, expectedCharacterJourney2, cache.characterJourneys["character2"])
	require.Equal(t, 1, len(cache.journeys))
	require.Equal(t, expectedJourney1, cache.journeys["23->42"])

	// only the first time a specific journey is started a NewJourney message is queue for DB
	mockQueue.AssertCalled(t, "Push", expectedMsg)
	mockQueue.AssertNumberOfCalls(t, "Push", 1)
}

func TestStartJourneySameShip(t *testing.T) {
	mockQueue := &queue.MockQueue{}
	cache := NewCache(nil, mockQueue)

	// first characterJourney of a character
	expectedMsg1 := queue.NewJourney{StartId: 23, DestinationId: 42}
	mockQueue.On("Push", expectedMsg1).Return()
	cache.StartJourney("character1", 23, 42)

	expectedCharacterJourney := &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedJourney1 := &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}
	require.Equal(t, 1, len(cache.characterJourneys))
	require.Equal(t, expectedCharacterJourney, cache.characterJourneys["character1"])
	require.Equal(t, 1, len(cache.journeys))
	require.Equal(t, expectedJourney1, cache.journeys["23->42"])

	// same character starts another route
	expectedMsg2 := queue.NewJourney{StartId: 42, DestinationId: 23}
	mockQueue.On("Push", expectedMsg2).Return()
	cache.StartJourney("character1", 42, 23)

	// just update entry (not creating new characterJourney) and create new journey
	expectedCharacterJourneyUpdated := &characterJourney{characterId: "character1", startId: 42, destinationId: 23}
	expectedJourney2 := &journey{
		startId: 42, destinationId: 23, points: nil, isFullyMapped: false,
	}
	require.Equal(t, 1, len(cache.characterJourneys))
	require.Equal(t, expectedCharacterJourneyUpdated, cache.characterJourneys["character1"])
	require.Equal(t, 2, len(cache.journeys))
	require.Equal(t, expectedJourney2, cache.journeys["42->23"])

	// both journeys are pushed into the queue
	mockQueue.AssertCalled(t, "Push", expectedMsg1)
	mockQueue.AssertCalled(t, "Push", expectedMsg2)
	mockQueue.AssertNumberOfCalls(t, "Push", 2)
}

func TestMovement(t *testing.T) {
	mockQueue := &queue.MockQueue{}
	cache := NewCache(nil, mockQueue)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	cache.characterJourneys["character2"] = &characterJourney{characterId: "character2", startId: 13, destinationId: 42}
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}
	cache.journeys["13->42"] = &journey{
		startId: 13, destinationId: 42, points: nil, isFullyMapped: false,
	}

	expectedMsg1 := queue.NewLocation{StartId: 23, DestinationId: 42, X: 1, Y: 2}
	mockQueue.On("Push", expectedMsg1).Return()
	err := cache.Movement("character1", 1, 2)
	require.NoError(t, err)

	// add new point to the correct journey
	require.Equal(t, []Point{{X: 1, Y: 2}}, cache.journeys["23->42"].points)
	require.Nil(t, cache.journeys["13->42"].points)

	expectedMsg2 := queue.NewLocation{StartId: 23, DestinationId: 42, X: 2, Y: 2}
	mockQueue.On("Push", expectedMsg2).Return()
	err = cache.Movement("character1", 2, 2)
	require.NoError(t, err)

	// append another point
	require.Equal(t, []Point{{X: 1, Y: 2}, {X: 2, Y: 2}}, cache.journeys["23->42"].points)
	require.Nil(t, cache.journeys["13->42"].points)

	// both journeys are pushed into the queue
	mockQueue.AssertCalled(t, "Push", expectedMsg1)
	mockQueue.AssertCalled(t, "Push", expectedMsg2)
	mockQueue.AssertNumberOfCalls(t, "Push", 2)
}

func TestMovementSamePoint(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedPoints := []Point{{X: 1, Y: 2}}
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: expectedPoints, isFullyMapped: false,
	}

	// ignore already existing point
	err := cache.Movement("character1", 1, 2)
	require.NoError(t, err)

	// no change
	require.Equal(t, expectedPoints, cache.journeys["23->42"].points)
}

func TestMovementErrorMissingVoyage(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedPoints := []Point{{X: 1, Y: 2}}
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: expectedPoints, isFullyMapped: false,
	}

	// no characterJourney for character2
	err := cache.Movement("character2", 11, 12)
	require.Error(t, err)
	require.Equal(t, "No active characterJourney for character character2", err.Error())

	// no change
	require.Equal(t, expectedPoints, cache.journeys["23->42"].points)
}

func TestMovementErrorMissingJourney(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedPoints := []Point{{X: 1, Y: 2}}
	cache.journeys["42->23"] = &journey{
		startId: 42, destinationId: 23, points: expectedPoints, isFullyMapped: false,
	}

	// no journey for character
	err := cache.Movement("character1", 11, 12)
	require.Error(t, err)
	require.Equal(t, "Missing journey between location 23 and 42", err.Error())

	// no change
	require.Equal(t, expectedPoints, cache.journeys["42->23"].points)
}

func TestReachedDestination(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockQueue := &queue.MockQueue{}
	cache := NewCache(mockMetrics, mockQueue)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	cache.characterJourneys["character2"] = &characterJourney{characterId: "character2", startId: 13, destinationId: 42}
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: []Point{{X: 1, Y: 2}}, isFullyMapped: false,
	}
	cache.journeys["13->42"] = &journey{
		startId: 13, destinationId: 42, points: []Point{{X: 11, Y: 12}}, isFullyMapped: false,
	}

	expectedMsg := queue.JourneyFullyMapped{StartId: 23, DestinationId: 42}
	mockQueue.On("Push", expectedMsg).Return()
	mockMetrics.On("LogJourney").Return()
	err := cache.ReachedDestination("character1", 42)
	require.NoError(t, err)

	// set the correct journey to 'isFullyMapped'
	require.True(t, cache.journeys["23->42"].isFullyMapped)
	require.False(t, cache.journeys["13->42"].isFullyMapped)

	// reach the end of the same journey again
	err = cache.ReachedDestination("character1", 42)
	require.NoError(t, err)

	// no change
	require.True(t, cache.journeys["23->42"].isFullyMapped)

	// only one msg is pushed into the queue
	mockQueue.AssertCalled(t, "Push", expectedMsg)
	mockQueue.AssertNumberOfCalls(t, "Push", 1)
	// route is only once counted
	mockMetrics.AssertNumberOfCalls(t, "LogJourney", 1)
}

func TestReachedDestinationErrorMissingVoyage(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedPoints := []Point{{X: 1, Y: 2}}
	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: expectedPoints, isFullyMapped: false,
	}

	// no characterJourney for character2
	err := cache.ReachedDestination("character2", 42)
	require.Error(t, err)
	require.Equal(t, "No active characterJourney for character character2", err.Error())

	// no change
	require.Equal(t, expectedPoints, cache.journeys["23->42"].points)
}

func TestReachedDestinationErrorMissingJourney(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.characterJourneys["character1"] = &characterJourney{characterId: "character1", startId: 23, destinationId: 42}
	expectedPoints := []Point{{X: 1, Y: 2}}
	cache.journeys["42->23"] = &journey{
		startId: 42, destinationId: 23, points: expectedPoints, isFullyMapped: false,
	}

	// no journey for character
	err := cache.ReachedDestination("character1", 42)
	require.Error(t, err)
	require.Equal(t, "Missing journey between location 23 and 42", err.Error())

	// no change
	require.Equal(t, expectedPoints, cache.journeys["42->23"].points)
}

func TestCheckJourneyNew(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}

	// new journey
	require.True(t, cache.checkJourney(13, 42))

	// cache is updated
	require.Equal(t, 2, len(cache.journeys))
	require.Equal(
		t,
		journey{startId: 13, destinationId: 42, points: nil, isFullyMapped: false},
		*cache.journeys["13->42"],
	)
}

func TestCheckJourneyExisting(t *testing.T) {
	cache := NewCache(nil, nil)

	cache.journeys["23->42"] = &journey{
		startId: 23, destinationId: 42, points: nil, isFullyMapped: false,
	}

	// existing journey
	require.False(t, cache.checkJourney(23, 42))

	// no change
	require.Equal(t, 1, len(cache.journeys))
}

func TestCheckPositionNew(t *testing.T) {
	route := &journey{
		startId: 23, destinationId: 42, points: []Point{{X: 1, Y: 2}}, isFullyMapped: false,
	}

	// new route point
	require.True(t, route.checkPosition(2, 2))

	// journey is updated
	require.Equal(t, 2, len(route.points))
	require.Equal(t, []Point{{X: 1, Y: 2}, {X: 2, Y: 2}}, route.points)
}

func TestCheckPositionExisting(t *testing.T) {
	route := &journey{
		startId: 23, destinationId: 42, points: []Point{{X: 1, Y: 2}}, isFullyMapped: false,
	}

	// existing route point
	require.False(t, route.checkPosition(1, 2))

	// no change
	require.Equal(t, 1, len(route.points))
	require.Equal(t, []Point{{X: 1, Y: 2}}, route.points)
}

func TestCheckPositionAlreadyMapped(t *testing.T) {
	route := &journey{
		startId: 23, destinationId: 42, points: []Point{{X: 1, Y: 2}}, isFullyMapped: true,
	}

	// point ignored
	require.False(t, route.checkPosition(2, 2))

	// no change
	require.Equal(t, 1, len(route.points))
	require.Equal(t, []Point{{X: 1, Y: 2}}, route.points)
}
