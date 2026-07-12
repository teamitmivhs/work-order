package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count int
	reset time.Time
}

// rateLimit applies a small in-memory fixed-window limit per client IP.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	clients := make(map[string]rateLimitEntry)

	return func(c *gin.Context) {
		now := time.Now()
		key := c.ClientIP()
		if forwarded := c.GetHeader("X-Real-IP"); net.ParseIP(forwarded) != nil {
			key = forwarded
		}

		mu.Lock()
		for client, item := range clients {
			if !now.Before(item.reset) {
				delete(clients, client)
			}
		}
		entry := clients[key]
		if entry.reset.IsZero() || !now.Before(entry.reset) {
			entry = rateLimitEntry{reset: now.Add(window)}
		}
		entry.count++
		clients[key] = entry
		mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		if entry.count > limit {
			c.Header("Retry-After", strconv.Itoa(max(1, int(time.Until(entry.reset).Seconds()))))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
