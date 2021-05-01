package queue

import (
	"errors"
	"time"
)

// TODO all operations should be fallible

var ERROR_TIMEOUT = errors.New("timeout")

type Queue interface {
	Empty() bool
	ReadQueue
	WriteQueue
}

type ReadQueue interface {
	BlockingRead() string
	ReadWithTimeout(time.Duration) (string, error)
}

type WriteQueue interface {
	Write(string)
	WriteMany([]string)
}
