package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/response"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

type IPRateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]*rateLimitEntry
}

func NewIPRateLimiter(limit int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		limit:   limit,
		window:  window,
		entries: make(map[string]*rateLimitEntry),
	}
}

func (r *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()
		ip := c.ClientIP()
		now := time.Now().UTC()

		r.mu.Lock()
		r.purgeExpiredLocked(now)

		entry, ok := r.entries[ip]
		if !ok || now.After(entry.resetAt) {
			entry = &rateLimitEntry{
				count:   0,
				resetAt: now.Add(r.window),
			}
			r.entries[ip] = entry
		}

		if entry.count >= r.limit {
			retryAfterSeconds := int(entry.resetAt.Sub(now).Seconds())
			if retryAfterSeconds < 0 {
				retryAfterSeconds = 0
			}
			log.Warn().
				Str("ip", ip).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Int("limit", r.limit).
				Int("retry_after_seconds", retryAfterSeconds).
				Msg("rate limit exceeded")
			r.mu.Unlock()
			response.FailureWithAbort(c, http.StatusTooManyRequests, "too many requests", "too many requests")
			return
		}

		entry.count++
		r.mu.Unlock()

		c.Next()
	}
}

func (r *IPRateLimiter) purgeExpiredLocked(now time.Time) {
	for ip, entry := range r.entries {
		if now.After(entry.resetAt) {
			delete(r.entries, ip)
		}
	}
}
