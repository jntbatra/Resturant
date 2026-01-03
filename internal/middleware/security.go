package middleware

import (
	"net/http"
	"sync"
	"time"

	apperrors "restaurant/internal/errors"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
	CleanupInterval   time.Duration
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	config   RateLimitConfig
	buckets  map[string]*TokenBucket
	mu       sync.RWMutex
	ticker   *time.Ticker
	stopChan chan struct{}
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens    float64
	maxTokens float64
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	if config.RequestsPerSecond == 0 {
		config.RequestsPerSecond = 100
	}
	if config.BurstSize == 0 {
		config.BurstSize = config.RequestsPerSecond * 2
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute
	}

	rl := &RateLimiter{
		config:   config,
		buckets:  make(map[string]*TokenBucket),
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	rl.ticker = time.NewTicker(config.CleanupInterval)
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given identifier is allowed
func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[identifier]
	if !exists {
		bucket = &TokenBucket{
			tokens:    float64(rl.config.BurstSize),
			maxTokens: float64(rl.config.BurstSize),
			lastReset: time.Now(),
		}
		rl.buckets[identifier] = bucket
	}

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.lastReset).Seconds()
	tokensToAdd := elapsed * float64(rl.config.RequestsPerSecond)
	bucket.tokens = min(bucket.tokens+tokensToAdd, bucket.maxTokens)
	bucket.lastReset = now

	// Check if we have tokens
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes old buckets periodically
func (rl *RateLimiter) cleanup() {
	for range rl.ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for id, bucket := range rl.buckets {
			if now.Sub(bucket.lastReset) > rl.config.CleanupInterval {
				delete(rl.buckets, id)
			}
		}
		rl.mu.Unlock()
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.stopChan)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Global rate limiter instance
var globalRateLimiter *RateLimiter

// InitRateLimiter initializes the global rate limiter
func InitRateLimiter(config RateLimitConfig) {
	globalRateLimiter = NewRateLimiter(config)
}

// RateLimitMiddleware is a Gin middleware for rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalRateLimiter == nil {
			c.Next()
			return
		}

		// Use IP address as identifier
		identifier := c.ClientIP()

		if !globalRateLimiter.Allow(identifier) {
			err := apperrors.NewAppError(
				http.StatusTooManyRequests,
				"rate limit exceeded",
				nil,
			)
			HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequestSizeLimitMiddleware limits the size of incoming requests
func RequestSizeLimitMiddleware(maxSizeBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSizeBytes)
		c.Next()
	}
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:5173", // Vite default
		"http://localhost:8080",
		"https://yourdomain.com",
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			// For public APIs, allow all but without credentials
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
