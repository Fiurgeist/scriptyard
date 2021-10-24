package server

import (
	"bytes"
	"fiurgeist/journey/internal/cache"
	"fiurgeist/journey/internal/metrics"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var defaultConfig = Config{
	Addr:        ":8080",
	PublicDir:   "foo",
	PublicJsDir: "foo/bar",
}

func TestMovementOK(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	mockMetrics.On("LogRequest").Return()
	mockCache.On("Movement", "character1", uint16(23), uint16(42)).Return(nil)

	jsonStr := []byte(`{"CharacterId": "character1", "X": 23, "Y": 42}`)
	req, _ := http.NewRequest("POST", "/character/movement", bytes.NewBuffer(jsonStr))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)
	mockMetrics.AssertCalled(t, "LogRequest")
	mockCache.AssertCalled(t, "Movement", "character1", uint16(23), uint16(42))
}

func TestMovementBadRequest(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	req, _ := http.NewRequest("POST", "/character/movement", bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusBadRequest, response.Code)
	mockMetrics.AssertNumberOfCalls(t, "LogRequest", 0)
	mockCache.AssertNumberOfCalls(t, "Movement", 0)
}

func TestReachedDestinationOK(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	mockMetrics.On("LogRequest").Return()
	mockCache.On("ReachedDestination", "character1", uint16(42)).Return(nil)

	jsonStr := []byte(`{"CharacterId": "character1", "DestinationId": 42}`)
	req, _ := http.NewRequest("POST", "/character/reachedDestination", bytes.NewBuffer(jsonStr))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)
	mockMetrics.AssertCalled(t, "LogRequest")
	mockCache.AssertCalled(t, "ReachedDestination", "character1", uint16(42))
}

func TestReachedDestinationBadRequest(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	req, _ := http.NewRequest("POST", "/character/reachedDestination", bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusBadRequest, response.Code)
	mockMetrics.AssertNumberOfCalls(t, "LogRequest", 0)
	mockCache.AssertNumberOfCalls(t, "ReachedDestination", 0)
}

func TestStartJourneyOK(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	mockMetrics.On("LogRequest").Return()
	mockCache.On("StartJourney", "character1", uint16(23), uint16(42)).Return(nil)

	jsonStr := []byte(`{"CharacterId": "character1", "StartId": 23, "DestinationId": 42, "X": 1, "Y": 2}`)
	req, _ := http.NewRequest("POST", "/character/startJourney", bytes.NewBuffer(jsonStr))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)
	mockMetrics.AssertCalled(t, "LogRequest")
	mockCache.AssertCalled(t, "StartJourney", "character1", uint16(23), uint16(42))
}

func TestStartJourneyBadRequest(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	req, _ := http.NewRequest("POST", "/character/startJourney", bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusBadRequest, response.Code)
	mockMetrics.AssertNumberOfCalls(t, "LogRequest", 0)
	mockCache.AssertNumberOfCalls(t, "StartJourney", 0)
}

func TestJourneysOK(t *testing.T) {
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(defaultConfig, mockMetrics, mockCache)

	journeyData := []cache.Journey{
		{Id: "23->42", Points: []cache.Point{{X: 1, Y: 2}, {X: 2, Y: 2}}},
		{Id: "42->23", Points: []cache.Point{{X: 11, Y: 12}, {X: 12, Y: 12}}},
	}
	mockCache.On("GetUniqueJourneys").Return(journeyData)

	req, _ := http.NewRequest("GET", "/journeys", bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)

	expected :=
		"{\"journeys\":[" +
			"{\"id\":\"23-\\u003e42\",\"data\":[{\"x\":1,\"y\":2},{\"x\":2,\"y\":2}]}," +
			"{\"id\":\"42-\\u003e23\",\"data\":[{\"x\":11,\"y\":12},{\"x\":12,\"y\":12}]}" +
			"]}\n"
	require.Equal(t, expected, response.Body.String())
}

func TestServeHome(t *testing.T) {
	publicDir, err := ioutil.TempDir("", "public")
	require.NoError(t, err)
	defer os.RemoveAll(publicDir)

	f, err := os.Create(filepath.Join(publicDir, "index.html"))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString("<html></html>")
	require.NoError(t, err)

	config := Config{
		Addr:        ":8080",
		PublicDir:   publicDir,
		PublicJsDir: "foo",
	}
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(config, mockMetrics, mockCache)

	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "<html></html>", response.Body.String())
}

func TestServeStaticFiles(t *testing.T) {
	publicJsDir, err := ioutil.TempDir("", "js")
	require.NoError(t, err)
	defer os.RemoveAll(publicJsDir)

	f, err := ioutil.TempFile(publicJsDir, "some_static.js")
	require.NoError(t, err)

	config := Config{
		Addr:        ":8080",
		PublicDir:   "foo",
		PublicJsDir: publicJsDir,
	}
	mockMetrics := &metrics.MockMetrics{}
	mockCache := &cache.MockCache{}
	srv := NewHTTPServer(config, mockMetrics, mockCache)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/static/js/%s", filepath.Base(f.Name())), bytes.NewBuffer([]byte("")))
	response := executeRequest(srv, req)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/javascript", response.Header()["Content-Type"][0])
}

func executeRequest(srv *http.Server, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	return rec
}
