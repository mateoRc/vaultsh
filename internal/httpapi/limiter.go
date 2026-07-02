package httpapi

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const maxRateLimitClients = 10000

type rateBucket struct {
	tokens float64
	last   time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]rateBucket
	now     func() time.Time
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]rateBucket),
		now:     time.Now,
	}
}

func (l *rateLimiter) allow(key string, perMinute, burst int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	bucket, found := l.buckets[key]
	if !found {
		if len(l.buckets) >= maxRateLimitClients {
			l.cleanup(now)
			if len(l.buckets) >= maxRateLimitClients {
				return false
			}
		}
		bucket = rateBucket{tokens: float64(burst), last: now}
	}

	elapsed := now.Sub(bucket.last).Seconds()
	bucket.tokens = min(
		float64(burst),
		bucket.tokens+elapsed*float64(perMinute)/60,
	)
	bucket.last = now
	if bucket.tokens < 1 {
		l.buckets[key] = bucket
		return false
	}
	bucket.tokens--
	l.buckets[key] = bucket
	return true
}

func (l *rateLimiter) cleanup(now time.Time) {
	for key, bucket := range l.buckets {
		if now.Sub(bucket.last) > 10*time.Minute {
			delete(l.buckets, key)
		}
	}
}

func clientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		forwarded := strings.TrimSpace(strings.Split(
			r.Header.Get("X-Forwarded-For"),
			",",
		)[0])
		if net.ParseIP(forwarded) != nil {
			return forwarded
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
