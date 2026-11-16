package ai

import (
	"sync"
	"time"
)

type toolRateLimiter struct {
	mu          sync.Mutex
	minInterval time.Duration
	last        map[string]time.Time
}

func newToolRateLimiter(minInterval time.Duration) *toolRateLimiter {
	return &toolRateLimiter{
		minInterval: minInterval,
		last:        make(map[string]time.Time),
	}
}

func (l *toolRateLimiter) Allow(tool string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if previous, ok := l.last[tool]; ok && now.Sub(previous) < l.minInterval {
		return false
	}

	l.last[tool] = now
	return true
}
