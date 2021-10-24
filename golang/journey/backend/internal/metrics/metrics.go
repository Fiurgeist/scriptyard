package metrics

import (
	"log"
	"sync/atomic"
	"time"
)

type Metrics interface {
	Close()
	LogRequest()
	LogJourney()
}

type metrics struct {
	requestCount uint64
	journeyCount uint64
	runningSince time.Time
	quit         chan bool
}

func NewMetrics() *metrics {
	quit := make(chan bool)
	m := newMetrics(quit)

	go func() {
		tick := time.Tick(time.Second)
		for {
			select {
			case <-m.quit:
				log.Println("Quit logging.")
				return
			case <-tick:
				m.print()
			}
		}
	}()
	return m
}

func newMetrics(quit chan bool) *metrics {
	return &metrics{
		requestCount: 0,
		journeyCount: 0,
		runningSince: time.Now(),
		quit:         quit,
	}
}

func (m *metrics) Close() {
	m.quit <- true
}

func (m *metrics) LogRequest() {
	atomic.AddUint64(&m.requestCount, 1)
}

func (m *metrics) LogJourney() {
	atomic.AddUint64(&m.journeyCount, 1)
}

func (m *metrics) print() {
	since := time.Since(m.runningSince)
	log.Printf(
		"Running since %s; received %.2f req/sec; %d unique journeys\n",
		since,
		float64(atomic.LoadUint64(&m.requestCount))/since.Seconds(),
		atomic.LoadUint64(&m.journeyCount),
	)
}
