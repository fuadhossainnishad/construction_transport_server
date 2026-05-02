package postgres

import (
	"math/rand"
	"time"
)

type RetryPolicy interface {
	Next() time.Duration
	HasNext() bool
	Reset()
}

type ExponentialBackoff struct {
	attempts   int
	maxRetries int
	base       time.Duration
	max        time.Duration
}

func NewExponentialBackoff(max int, base, maxDelay time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		maxRetries: max,
		base:       base,
		max:        maxDelay,
	}
}

func (b *ExponentialBackoff) HasNext() bool {
	return b.attempts < b.maxRetries
}

func (b *ExponentialBackoff) Reset() {
	b.attempts = 0
}

func (b *ExponentialBackoff) Next() time.Duration {
	b.attempts++

	delay := time.Duration(1<<b.attempts) * b.base

	if delay > b.max {
		delay = b.max
	}

	jitter := time.Duration(rand.Intn(300)) * time.Millisecond
	return delay + jitter
}
