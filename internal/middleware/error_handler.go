package middleware

import (
	"fmt"
	"net/http"

	"restaurant/internal/errors"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the standardized error response format
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// ErrorHandler is a Gin middleware that handles errors from handlers
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if any errors occurred
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			// Try to extract AppError
			var statusCode int
			var errorMessage string

			if appErr, ok := err.Err.(*errors.AppError); ok {
				statusCode = appErr.StatusCode()
				errorMessage = appErr.Message
				if appErr.Err != nil {
					// Include underlying error details in development
					errorMessage = fmt.Sprintf("%s: %v", appErr.Message, appErr.Err)
				}
			} else {
				// Default to 500 for unknown errors
				statusCode = http.StatusInternalServerError
				errorMessage = "internal server error"
			}

			// Write error response if not already written
			if !c.Writer.Written() {
				c.JSON(statusCode, ErrorResponse{
					Error:  errorMessage,
					Status: statusCode,
				})
			}
		}
	}
}

// HandleError is a helper function to properly add errors to Gin context
// Usage: if err != nil { middleware.HandleError(c, err); return }
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
}

// HandleErrorWithCode wraps error with custom AppError type
// Usage: if err != nil { middleware.HandleErrorWithCode(c, http.StatusBadRequest, "custom message", err); return }
func HandleErrorWithCode(c *gin.Context, code int, message string, err error) {
	appErr := errors.NewAppError(code, message, err)
	c.Error(appErr)
}
