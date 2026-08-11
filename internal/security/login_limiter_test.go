package security

import (
	"testing"
	"time"
)

func TestLoginLimiter(t *testing.T) {
	limiter := NewLoginLimiter(2, time.Minute)
	now := time.Unix(1_000, 0)
	if !limiter.Allow("client", now) {
		t.Fatal("new client should be allowed")
	}
	limiter.RecordFailure("client", now)
	limiter.RecordFailure("client", now)
	if limiter.Allow("client", now) {
		t.Fatal("client should be blocked after reaching the limit")
	}
	if !limiter.Allow("client", now.Add(time.Minute)) {
		t.Fatal("client should be allowed after the window")
	}
}

func TestLoginLimiterReset(t *testing.T) {
	limiter := NewLoginLimiter(1, time.Minute)
	now := time.Unix(1_000, 0)
	limiter.RecordFailure("client", now)
	limiter.Reset("client")
	if !limiter.Allow("client", now) {
		t.Fatal("successful login should reset failures")
	}
}
