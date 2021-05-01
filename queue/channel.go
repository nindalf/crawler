package queue

import "time"

type channelQueue struct {
	ch chan string
}

func NewChannelQueue(size int) Queue {
	ch := make(chan string, size)
	return channelQueue{ch}
}

func (q channelQueue) Empty() bool {
	return len(q.ch) == 0
}

func (q channelQueue) BlockingRead() string {
	return <-q.ch
}

func (q channelQueue) ReadWithTimeout(timeout time.Duration) (string, error) {
	select {
	case val := <-q.ch:
		return val, nil
	case <-time.After(timeout):
		return "", ErrTimeout
	}
}

func (q channelQueue) Write(msg string) {
	q.ch <- msg
}

func (q channelQueue) WriteMany(msgs []string) {
	for _, msg := range msgs {
		q.ch <- msg
	}
}
