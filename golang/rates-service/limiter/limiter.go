package limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	lock   sync.Mutex
	tokens int
}

func NewRateLimiter() *RateLimiter {
	rateLimiter := &RateLimiter{}
	go rateLimiter.Refill()
	return rateLimiter
}

func (r *RateLimiter) Allow() bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.tokens > 0 {
		r.tokens -= 1
		return true
	}
	return false
}

func (r *RateLimiter) Refill() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		r.lock.Lock()
		r.tokens = 1
		r.lock.Unlock()
	}
}
