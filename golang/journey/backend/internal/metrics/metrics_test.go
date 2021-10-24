package metrics

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/undefinedlabs/go-mpatch"
	"log"
	"os"
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	quit := make(chan bool)
	metrics := newMetrics(quit)

	go metrics.Close()
	assert.Eventually(t, func() bool { return <-quit }, time.Second, 10*time.Millisecond)
}

func TestLogRequest(t *testing.T) {
	metrics := newMetrics(nil)

	require.Equal(t, uint64(0), metrics.requestCount)
	metrics.LogRequest()
	require.Equal(t, uint64(1), metrics.requestCount)
}

func TestLogJourney(t *testing.T) {
	metrics := newMetrics(nil)

	require.Equal(t, uint64(0), metrics.journeyCount)
	metrics.LogJourney()
	require.Equal(t, uint64(1), metrics.journeyCount)
}

func TestPrint(t *testing.T) {
	patchRunningSince, err := mpatch.PatchMethod(time.Now, mockRunningSince)
	require.NoError(t, err)
	defer patchRunningSince.Unpatch()

	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	metrics := newMetrics(nil)

	// test print with default values
	metrics.print()
	require.True(t, scanner.Scan())
	require.Equal(
		t,
		"2021/01/01 00:00:00 Running since 0s; received NaN req/sec; 0 unique journeys",
		scanner.Text(),
	)

	err = patchRunningSince.Unpatch()
	require.NoError(t, err)

	// with updated values
	metrics.requestCount = uint64(42)
	metrics.journeyCount = uint64(23)

	// test print two seconds later
	patchNow, err := mpatch.PatchMethod(time.Now, mockNow)
	require.NoError(t, err)
	defer patchNow.Unpatch()

	metrics.print()
	require.True(t, scanner.Scan())
	require.Equal(
		t,
		"2021/01/01 00:00:02 Running since 2s; received 21.00 req/sec; 23 unique journeys",
		scanner.Text(),
	)
}

func TestNewMetrics(t *testing.T) {
	patchRunningSince, err := mpatch.PatchMethod(time.Now, mockRunningSince)
	require.NoError(t, err)
	defer patchRunningSince.Unpatch()

	scanner, reader, writer := mockLogger(t)
	defer resetLogger(reader, writer)

	metrics := NewMetrics()

	// test printing
	require.True(t, scanner.Scan())
	require.Equal(
		t,
		"2021/01/01 00:00:00 Running since 0s; received NaN req/sec; 0 unique journeys",
		scanner.Text(),
	)

	// test quiting subroutine
	metrics.Close()
	require.True(t, scanner.Scan())
	require.Equal(t, "2021/01/01 00:00:00 Quit logging.", scanner.Text())
}

func mockRunningSince() time.Time {
	return time.Date(2021, 01, 01, 00, 00, 00, 0, time.UTC)
}

func mockNow() time.Time {
	return time.Date(2021, 01, 01, 00, 00, 02, 0, time.UTC)
}

// https://stackoverflow.com/a/68807754
func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetLogger(reader *os.File, writer *os.File) {
	defer log.SetOutput(os.Stderr)
	err := reader.Close()
	if err != nil {
		fmt.Println("Error closing reader: ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("Error closing writer: ", err)
	}
}
