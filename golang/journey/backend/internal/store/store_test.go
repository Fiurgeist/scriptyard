package store

import (
	"database/sql"
	"fiurgeist/journey/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) { // TODO: split this test
	quitSubroutine := false
	mockQueue := &queue.MockQueue{}
	store, err := NewStore(mockQueue)
	require.NoError(t, err)
	defer func() {
		if quitSubroutine {
			store.db.Close()
		} else {
			store.Close()
		}
	}()

	channel := make(chan interface{})
	mockQueue.On("GetChannel").Return(channel)

	// assert empty tables
	rows, err := store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	assertJourneyRows(t, rows, []journeyRow{})
	rows, err = store.db.Query("SELECT * FROM location WHERE 1;")
	require.NoError(t, err)
	assertLocationRows(t, rows, []locationRow{})

	// assert reading the three msg types from queue
	channel <- queue.NewJourney{StartId: 23, DestinationId: 42}
	assert.Eventually(t, func() bool { return len(channel) == 0 }, time.Second, 10*time.Millisecond)
	rows, err = store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	assertJourneyRows(
		t,
		rows,
		[]journeyRow{{id: "23-42", start: 23, end: 42, fullyMapped: false}},
	)

	channel <- queue.NewLocation{StartId: 23, DestinationId: 42, X: 1, Y: 2}
	assert.Eventually(t, func() bool { return len(channel) == 0 }, time.Second, 10*time.Millisecond)
	rows, err = store.db.Query("SELECT * FROM location WHERE 1;")
	require.NoError(t, err)
	assertLocationRows(
		t,
		rows,
		[]locationRow{{journeyId: "23-42", x: 1, y: 2, hackyUnique: "23-42-1-2"}},
	)

	channel <- queue.JourneyFullyMapped{StartId: 23, DestinationId: 42}
	assert.Eventually(t, func() bool { return len(channel) == 0 }, time.Second, 10*time.Millisecond)
	rows, err = store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	expectedJourneys := []journeyRow{{id: "23-42", start: 23, end: 42, fullyMapped: true}}
	assertJourneyRows(t, rows, expectedJourneys)

	// assert subroutine is cleaned up
	quitSubroutine = true
	store.subroutineQuit <- true
	store.subroutineWG.Wait()

	// TODO: find a better way to test this, without using basically "sleep"
	unconsumedMsg := queue.NewJourney{StartId: 42, DestinationId: 23}
	timeout := time.After(time.Second)
	go func() {
		channel <- unconsumedMsg
	}()
	<-timeout // timeout expected
	require.Equal(t, unconsumedMsg, <-channel)

	// no new journey
	rows, err = store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	assertJourneyRows(t, rows, expectedJourneys)
}

func TestClose(t *testing.T) {
	dbClosed := uint32(0)
	db, err := sql.Open("ramsql", "TestClose")
	require.NoError(t, err)
	defer func() {
		if atomic.LoadUint32(&dbClosed) == 0 {
			db.Close()
		}
	}()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	// simulate subroutine
	store.subroutineWG.Add(1)
	// assert DB open
	require.Equal(t, nil, store.db.Ping())

	go func() {
		store.Close()
		atomic.StoreUint32(&dbClosed, 1)
	}()

	go func() {
		// simulate subroutine
		<-store.subroutineQuit
		store.subroutineWG.Done()
	}()

	assert.Eventually(
		t,
		func() bool { return atomic.LoadUint32(&dbClosed) == 1 },
		time.Second,
		10*time.Millisecond,
	)

	// assert DB closed
	err = store.db.Ping()
	require.Error(t, err)
	require.Equal(t, "sql: database is closed", err.Error())
}

func TestNewJourney(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewJourney")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	err = store.newJourney(23, 42)
	require.NoError(t, err)

	rows, err := store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	journeyData := []journeyRow{{id: "23-42", start: 23, end: 42, fullyMapped: false}}
	assertJourneyRows(t, rows, journeyData)
}

func TestNewJourneyHandleDuplicate(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewJourneyHandleDuplicate")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	err = store.newJourney(23, 42)
	require.NoError(t, err)

	rows, err := store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	journeyData := []journeyRow{{id: "23-42", start: 23, end: 42, fullyMapped: false}}
	assertJourneyRows(t, rows, journeyData)

	// no error but same data
	err = store.newJourney(23, 42)
	require.NoError(t, err)

	rows, err = store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	assertJourneyRows(t, rows, journeyData)
}

func TestNewJourneyErrorInsertJourney(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewJourneyErrorInsertJourney")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)

	err = store.newJourney(23, 42)
	require.Error(t, err)
	require.Equal(t, "table journey does not exists", err.Error())
}

func TestJourneyFullyMapped(t *testing.T) {
	db, err := sql.Open("ramsql", "TestJourneyFullyMapped")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	// add some journeys
	err = store.newJourney(23, 42)
	require.NoError(t, err)
	err = store.newJourney(42, 23)
	require.NoError(t, err)

	// update one journey
	err = store.journeyFullyMapped(23, 42)
	require.NoError(t, err)

	rows, err := store.db.Query("SELECT * FROM journey WHERE 1;")
	require.NoError(t, err)
	journeyData := []journeyRow{
		{id: "23-42", start: 23, end: 42, fullyMapped: true}, // only one journey is changed
		{id: "42-23", start: 42, end: 23, fullyMapped: false},
	}
	assertJourneyRows(t, rows, journeyData)
}

func TestTestJourneyFullyMappedError(t *testing.T) {
	db, err := sql.Open("ramsql", "TestTestJourneyFullyMappedError")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)

	err = store.journeyFullyMapped(23, 42)
	require.Error(t, err)
	require.Equal(t, "Table journey does not exists", err.Error())
}

func TestNewLocation(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewLocation")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	err = store.newLocation(23, 42, 1, 2)
	require.NoError(t, err)

	// several trade journeys can share the same coordinate
	err = store.newLocation(42, 23, 1, 2)
	require.NoError(t, err)

	rows, err := store.db.Query("SELECT * FROM location WHERE 1;")
	require.NoError(t, err)
	locationData := []locationRow{
		{journeyId: "23-42", x: 1, y: 2, hackyUnique: "23-42-1-2"},
		{journeyId: "42-23", x: 1, y: 2, hackyUnique: "42-23-1-2"},
	}
	assertLocationRows(t, rows, locationData)
}

func TestNewLocationHandleDuplicate(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewLocationHandleDuplicate")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)
	err = store.init()
	require.NoError(t, err)

	err = store.newLocation(23, 42, 1, 2)
	require.NoError(t, err)

	rows, err := store.db.Query("SELECT * FROM location WHERE 1;")
	require.NoError(t, err)
	locationData := []locationRow{{journeyId: "23-42", x: 1, y: 2, hackyUnique: "23-42-1-2"}}
	assertLocationRows(t, rows, locationData)

	// no error but same data
	err = store.newLocation(23, 42, 1, 2)
	require.NoError(t, err)

	rows, err = store.db.Query("SELECT * FROM location WHERE 1;")
	require.NoError(t, err)
	assertLocationRows(t, rows, locationData)
}

func TestNewLocationError(t *testing.T) {
	db, err := sql.Open("ramsql", "TestNewLocationError")
	require.NoError(t, err)
	defer db.Close()

	store := newStore(db, nil)

	err = store.newLocation(23, 42, 1, 2)
	require.Error(t, err)
	require.Equal(t, "table location does not exists", err.Error())
}

type journeyRow struct {
	id          string
	start       uint16
	end         uint16
	fullyMapped bool
}

type locationRow struct {
	journeyId   string
	x           uint16
	y           uint16
	hackyUnique string
}

func assertJourneyRows(t *testing.T, rows *sql.Rows, expected []journeyRow) {
	nb := 0
	for rows.Next() {
		var gotId string
		var gotStart, gotDest uint16
		var gotFullyMapped bool
		err := rows.Scan(&gotId, &gotStart, &gotDest, &gotFullyMapped)
		require.NoError(t, err)
		require.LessOrEqual(t, nb, len(expected))
		require.Equal(t, expected[nb].id, gotId)
		require.Equal(t, expected[nb].start, gotStart)
		require.Equal(t, expected[nb].end, gotDest)
		require.Equal(t, expected[nb].fullyMapped, gotFullyMapped)
		nb++
	}
	require.Equal(t, len(expected), nb)
}

func assertLocationRows(t *testing.T, rows *sql.Rows, expected []locationRow) {
	nb := 0
	for rows.Next() {
		var gotJourneyId, gotHackyUnique string
		var gotX, gotY uint16
		err := rows.Scan(&gotJourneyId, &gotX, &gotY, &gotHackyUnique)
		require.NoError(t, err)
		require.LessOrEqual(t, nb, len(expected))
		require.Equal(t, expected[nb].journeyId, gotJourneyId)
		require.Equal(t, expected[nb].x, gotX)
		require.Equal(t, expected[nb].y, gotY)
		require.Equal(t, expected[nb].hackyUnique, gotHackyUnique)
		nb++
	}
	require.Equal(t, len(expected), nb)
}
