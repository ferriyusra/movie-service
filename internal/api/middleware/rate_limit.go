package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	ips      map[string]*ipLimiter
	r        rate.Limit
	burst    int
}

func newRateLimiter(r rate.Limit, burst int) *rateLimiter {
	rl := &rateLimiter{
		ips:   make(map[string]*ipLimiter),
		r:     r,
		burst: burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, ok := rl.ips[ip]
	if !ok {
		entry = &ipLimiter{limiter: rate.NewLimiter(rl.r, rl.burst)}
		rl.ips[ip] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

// cleanup removes IPs not seen in the last 5 minutes to prevent memory leaks
func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		for ip, entry := range rl.ips {
			if time.Since(entry.lastSeen) > 5*time.Minute {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit returns a middleware that limits requests to `limit` per minute per IP
func RateLimit(perMinute int) gin.HandlerFunc {
	rl := newRateLimiter(rate.Every(time.Minute/time.Duration(perMinute)), perMinute)

	return func(c *gin.Context) {
		if !rl.get(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Err("Too many requests, please try again later"))
			return
		}
		c.Next()
	}
}
