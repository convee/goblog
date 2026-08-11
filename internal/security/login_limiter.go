package security

import (
	"sync"
	"time"
)

type loginAttempt struct {
	failures int
	resetAt  time.Time
}

type LoginLimiter struct {
	mu       sync.Mutex
	attempts map[string]loginAttempt
	limit    int
	window   time.Duration
}

func NewLoginLimiter(limit int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{
		attempts: make(map[string]loginAttempt),
		limit:    limit,
		window:   window,
	}
}

func (l *LoginLimiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	attempt, exists := l.attempts[key]
	if !exists || !now.Before(attempt.resetAt) {
		delete(l.attempts, key)
		return true
	}
	return attempt.failures < l.limit
}

func (l *LoginLimiter) RecordFailure(key string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	attempt, exists := l.attempts[key]
	if !exists || !now.Before(attempt.resetAt) {
		attempt = loginAttempt{resetAt: now.Add(l.window)}
	}
	attempt.failures++
	l.attempts[key] = attempt
	if len(l.attempts) > 10_000 {
		l.removeExpired(now)
	}
}

func (l *LoginLimiter) Reset(key string) {
	l.mu.Lock()
	delete(l.attempts, key)
	l.mu.Unlock()
}

func (l *LoginLimiter) removeExpired(now time.Time) {
	for key, attempt := range l.attempts {
		if !now.Before(attempt.resetAt) {
			delete(l.attempts, key)
		}
	}
}
