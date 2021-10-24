package server

import (
	"encoding/json"
	"fiurgeist/journey/internal/cache"
	"fiurgeist/journey/internal/metrics"
	"github.com/gorilla/mux"
	"net/http"
)

type Config struct {
	Addr        string
	PublicDir   string
	PublicJsDir string
}

func NewHTTPServer(config Config, metrics metrics.Metrics, cache cache.Cache) *http.Server {
	httpsrv := newHTTPServer(metrics, cache)
	r := mux.NewRouter()

	r.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.PublicDir))))
	r.PathPrefix("/static/js").HandlerFunc(
		setHeader(
			"Content-Type",
			"application/javascript",
			http.StripPrefix("/static/js", http.FileServer(http.Dir(config.PublicJsDir))),
		),
	)

	r.HandleFunc("/character/movement", httpsrv.handleMovement).Methods("POST")
	r.HandleFunc("/character/startJourney", httpsrv.handleStartJourney).Methods("POST")
	r.HandleFunc("/character/reachedDestination", httpsrv.handleReachedDestination).Methods("POST")

	r.HandleFunc("/journeys", httpsrv.handleJourneys).Methods("GET", "OPTIONS")

	return &http.Server{
		Addr:    config.Addr,
		Handler: r,
	}
}

func setHeader(header, value string, handle http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(header, value)
		handle.ServeHTTP(w, req)
	}
}

type httpServer struct {
	metrics metrics.Metrics
	cache   cache.Cache
}

func newHTTPServer(metrics metrics.Metrics, cache cache.Cache) *httpServer {
	return &httpServer{
		metrics: metrics,
		cache:   cache,
	}
}

type MovementRequest struct {
	CharacterId string `json:"CharacterId"`
	X           uint16 `json:"X"`
	Y           uint16 `json:"Y"`
}

type ReachedDestinationRequest struct {
	CharacterId   string `json:"CharacterId"`
	DestinationId uint16 `json:"DestinationId"`
}

type StartJourneyRequest struct {
	CharacterId   string `json:"CharacterId"`
	StartId       uint16 `json:"StartId"`
	DestinationId uint16 `json:"DestinationId"`
}

type JourneysResponse struct {
	Journeys []cache.Journey `json:"journeys"`
}

func (s *httpServer) handleMovement(w http.ResponseWriter, r *http.Request) {
	var req MovementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metrics.LogRequest()
	s.cache.Movement(req.CharacterId, req.X, req.Y)
}

func (s *httpServer) handleReachedDestination(w http.ResponseWriter, r *http.Request) {
	var req ReachedDestinationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metrics.LogRequest()
	s.cache.ReachedDestination(req.CharacterId, req.DestinationId)
}

func (s *httpServer) handleStartJourney(w http.ResponseWriter, r *http.Request) {
	var req StartJourneyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metrics.LogRequest()
	s.cache.StartJourney(req.CharacterId, req.StartId, req.DestinationId)
}

func (s *httpServer) handleJourneys(w http.ResponseWriter, r *http.Request) {
	res := JourneysResponse{
		Journeys: s.cache.GetUniqueJourneys(),
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
