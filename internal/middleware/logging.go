package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the context key for storing request ID
const RequestIDKey = "request-id"

// LogEntry represents a structured log entry
type LogEntry struct {
	RequestID string
	Timestamp string
	Method    string
	Path      string
	Status    int
	Duration  string
	Message   string
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Get request info
		method := c.Request.Method
		path := c.Request.URL.Path
		requestID := c.GetString(RequestIDKey)

		// Process request
		c.Next()

		// Get response info
		statusCode := c.Writer.Status()
		duration := time.Since(startTime)

		// Log request with structured data
		entry := LogEntry{
			RequestID: requestID,
			Timestamp: startTime.Format(time.RFC3339),
			Method:    method,
			Path:      path,
			Status:    statusCode,
			Duration:  duration.String(),
		}

		// Determine log level based on status code
		if statusCode >= 400 {
			entry.Message = "request failed"
		} else {
			entry.Message = "request success"
		}

		log.Printf("[%s] %s %s %s %d %s (%s)",
			entry.RequestID,
			entry.Timestamp,
			entry.Method,
			entry.Path,
			entry.Status,
			entry.Message,
			entry.Duration,
		)
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	requestID := c.GetString(RequestIDKey)
	if requestID == "" {
		requestID = "unknown"
	}
	return requestID
}

// LogOperation logs a database or service operation with request ID
func LogOperation(c *gin.Context, operation string, details string) {
	requestID := GetRequestID(c)
	log.Printf("[%s] %s: %s", requestID, operation, details)
}
