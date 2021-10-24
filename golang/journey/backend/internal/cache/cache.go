package cache

import (
	"fiurgeist/journey/internal/metrics"
	"fiurgeist/journey/internal/queue"
	"fmt"
	"log"
	"sync"
)

type Journey struct {
	Id     string  `json:"id"`
	Points []Point `json:"data"`
}

type Point struct {
	X uint16 `json:"x"`
	Y uint16 `json:"y"`
}

type Cache interface {
	GetUniqueJourneys() []Journey
	StartJourney(characterId string, startId, destinationId uint16)
	Movement(characterId string, x, y uint16) error
	ReachedDestination(characterId string, destinationId uint16) error
}

type characterJourney struct {
	characterId   string
	startId       uint16
	destinationId uint16
}

type journey struct {
	startId          uint16
	destinationId    uint16
	points           []Point
	isFullyMapped    bool
}

type cache struct {
	mu                sync.RWMutex
	characterJourneys map[string]*characterJourney
	journeys          map[string]*journey
	msgQueue          queue.Queue
	metrics           metrics.Metrics
}

func NewCache(metrics metrics.Metrics, msgQueue queue.Queue) *cache {
	r := &cache{
		characterJourneys: make(map[string]*characterJourney),
		journeys:          make(map[string]*journey),
		msgQueue:          msgQueue,
		metrics:           metrics,
	}
	return r
}

func (c *cache) GetUniqueJourneys() []Journey {
	c.mu.RLock()
	defer c.mu.RUnlock()

	routes := make([]Journey, len(c.journeys))
	index := 0
	for id, route := range c.journeys {
		routes[index] = Journey{Id: id, Points: route.points}
		index++
	}
	return routes
}

func (c *cache) StartJourney(characterId string, startId, destinationId uint16) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.characterJourneys[characterId] = &characterJourney{
		characterId:   characterId,
		startId:       startId,
		destinationId: destinationId,
	}
	if c.checkJourney(startId, destinationId) {
		c.msgQueue.Push(queue.NewJourney{
			StartId:       startId,
			DestinationId: destinationId,
		})
	}
}

func (c *cache) Movement(characterId string, x, y uint16) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	characterJourney := c.characterJourneys[characterId]
	if characterJourney == nil {
		err := fmt.Errorf("No active characterJourney for character %s", characterId)
		log.Println(err.Error())
		return err
	}

	routeKey := fmt.Sprintf("%d->%d", characterJourney.startId, characterJourney.destinationId)
	route := c.journeys[routeKey]
	if route == nil {
		err := fmt.Errorf(
			"Missing journey between location %d and %d", characterJourney.startId, characterJourney.destinationId,
		)
		log.Println(err.Error())
		return err
	}
	if route.checkPosition(x, y) {
		c.msgQueue.Push(queue.NewLocation{
			StartId:       characterJourney.startId,
			DestinationId: characterJourney.destinationId,
			X:             x,
			Y:             y,
		})
	}

	return nil
}

func (c *cache) ReachedDestination(characterId string, destinationId uint16) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	characterJourney := c.characterJourneys[characterId]
	if characterJourney == nil {
		err := fmt.Errorf("No active characterJourney for character %s", characterId)
		log.Println(err.Error())
		return err
	}

	routeKey := fmt.Sprintf("%d->%d", characterJourney.startId, characterJourney.destinationId)
	route := c.journeys[routeKey]
	if route == nil {
		err := fmt.Errorf(
			"Missing journey between location %d and %d", characterJourney.startId, characterJourney.destinationId,
		)
		log.Println(err.Error())
		return err
	}

	if route.isFullyMapped {
		return nil
	}

	route.isFullyMapped = true
	c.metrics.LogJourney()
	c.msgQueue.Push(queue.JourneyFullyMapped{
		StartId:       characterJourney.startId,
		DestinationId: characterJourney.destinationId,
	})

	return nil
}

func (c *cache) checkJourney(startId, destinationId uint16) bool {
	routeKey := fmt.Sprintf("%d->%d", startId, destinationId)
	route := c.journeys[routeKey]
	if route == nil {
		c.journeys[routeKey] = &journey{
			startId:       startId,
			destinationId: destinationId,
			isFullyMapped: false,
		}
		return true
	}
	return false
}

func (t *journey) checkPosition(x, y uint16) bool {
	if t.isFullyMapped {
		return false
	}
	for _, point := range t.points {
		if point.X == x && point.Y == y {
			return false
		}
	}
	t.points = append(t.points, Point{X: x, Y: y})
	return true
}
