package queue

const CHANNEL_BUFFER_SIZE = 1024 * 1024

type NewJourney struct {
	StartId       uint16
	DestinationId uint16
}
type NewLocation struct {
	StartId       uint16
	DestinationId uint16
	X             uint16
	Y             uint16
}
type JourneyFullyMapped struct {
	StartId       uint16
	DestinationId uint16
}

type Queue interface {
	Close()
	Push(msg interface{})
	GetChannel() chan interface{}
}

type queue struct {
	channel chan interface{}
}

func NewQueue() *queue {
	q := &queue{
		channel: make(chan interface{}, CHANNEL_BUFFER_SIZE),
	}
	return q
}

func (q *queue) Close() {
	close(q.channel)
}

func (q *queue) Push(msg interface{}) {
	q.channel <- msg
}

func (q *queue) GetChannel() chan interface{} {
	return q.channel
}
