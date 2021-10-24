package store

import (
	"database/sql"
	"fiurgeist/journey/internal/queue"
	"fmt"
	_ "github.com/proullon/ramsql/driver"
	"log"
	"sync"
)

type store struct {
	db             *sql.DB
	msgQueue       queue.Queue
	subroutineQuit chan bool
	subroutineWG   *sync.WaitGroup
}

func NewStore(msgQueue queue.Queue) (*store, error) {
	db, err := sql.Open("ramsql", "JourneyDB")
	if err != nil {
		fmt.Printf("sql.Open : Error : %s\n", err)
		return nil, err
	}

	s := newStore(db, msgQueue)
	if err := s.init(); err != nil {
		return nil, err
	}

	s.subroutineWG.Add(1)
	go func() {
		defer s.subroutineWG.Done()
		for {
			select {
			case <-s.subroutineQuit:
				log.Println("Quit storing.")
				return
			case msg := <-s.msgQueue.GetChannel():
				switch data := msg.(type) {
				case queue.NewJourney:
					s.newJourney(data.StartId, data.DestinationId)
				case queue.JourneyFullyMapped:
					s.journeyFullyMapped(data.StartId, data.DestinationId)
				case queue.NewLocation:
					s.newLocation(data.StartId, data.DestinationId, data.X, data.Y)
				}
			}
		}
	}()
	return s, nil
}

func newStore(db *sql.DB, msgQueue queue.Queue) *store {
	return &store{
		db:             db,
		msgQueue:       msgQueue,
		subroutineQuit: make(chan bool),
		subroutineWG:   &sync.WaitGroup{},
	}
}

func (s *store) init() error {
	batch := []string{
		`CREATE TABLE journey (id TEXT UNIQUE NOT NULL, start_id INT , destination_id INT, fully_mapped BOOLEAN);`,
		`CREATE TABLE location (journey_id INT, x INT, y INT, ramsql_hack_unique_composite_key TEXT UNIQUE NOT NULL);`,
	}

	for _, b := range batch {
		_, err := s.db.Exec(b)
		if err != nil {
			log.Printf("sql.Exec: Error: %s\n", err)
			return err
		}
	}

	return nil
}

func (s *store) Close() {
	s.subroutineQuit <- true
	s.subroutineWG.Wait()
	s.db.Close()
}

func (s *store) newJourney(startId, destinationId uint16) error {
	query := `INSERT INTO journey (id, start_id, destination_id, fully_mapped) VALUES ($1, $2, $3, 'FALSE');`
	_, err := s.db.Exec(
		query, fmt.Sprintf("%d-%d", startId, destinationId), startId, destinationId,
	)
	if err != nil && err.Error() != "UNIQUE constraint violation" {
		log.Printf(
			"Failed to insert new journey: %s; (startId: %d, destinationId: %d)\n",
			err,
			startId,
			destinationId,
		)
		return err
	}

	return nil
}

func (s *store) journeyFullyMapped(startId, destinationId uint16) error {
	query := `UPDATE journey SET fully_mapped = 'TRUE' WHERE start_id = $1 AND destination_id = $2;`
	_, err := s.db.Exec(query, startId, destinationId)
	if err != nil {
		log.Printf(
			"Failed to update `fully_mapped` of journey : %s; (startId: %d, destinationId: %d)\n",
			err,
			startId,
			destinationId,
		)
		return err
	}
	return nil
}

func (s *store) newLocation(startId, destinationId, x, y uint16) error {
	query := `INSERT INTO location (journey_id, x, y, ramsql_hack_unique_composite_key)
                VALUES ($1, $2, $3, $4);`
	journey_id := fmt.Sprintf("%d-%d", startId, destinationId)
	_, err := s.db.Exec(query, journey_id, x, y, fmt.Sprintf("%s-%d-%d", journey_id, x, y))
	if err != nil && err.Error() != "UNIQUE constraint violation" {
		log.Printf(
			"Failed to insert new location: %s; (startId: %d, destinationId: %d, x: %d, y: %d)\n",
			err,
			startId,
			destinationId,
			x,
			y,
		)
		return err
	}
	return nil
}
