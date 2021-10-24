package queue

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewQueue(t *testing.T) {
	queue := NewQueue()

	require.Equal(t, 0, len(queue.channel))
	require.Equal(t, CHANNEL_BUFFER_SIZE, cap(queue.channel))
}

func TestClose(t *testing.T) {
	queue := NewQueue()

	queue.Close()
	_, ok := <-queue.channel
	require.False(t, ok)
}

func Push(t *testing.T) {
	queue := NewQueue()

	require.Equal(t, 0, len(queue.channel))

	queue.Push(struct{}{})

	require.Equal(t, 1, len(queue.channel))
}

func GetChannel(t *testing.T) {
	queue := NewQueue()
	queue.Push(struct{ foo int }{foo: 42})

	require.Equal(t, 1, len(queue.GetChannel()))
	require.Equal(t, struct{ foo int }{foo: 42}, <-queue.GetChannel())
}
