package middleware

import (
	apperrors "restaurant/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UUIDParam extracts and parses a UUID parameter from the request
// Returns an error and writes an error response if parsing fails
func UUIDParam(c *gin.Context, paramName string) (uuid.UUID, bool) {
	paramStr := c.Param(paramName)
	if paramStr == "" {
		err := apperrors.NewValidationError("missing " + paramName + " parameter")
		HandleError(c, err)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(paramStr)
	if err != nil {
		err := apperrors.NewValidationError("invalid " + paramName + ": must be a valid UUID")
		HandleError(c, err)
		return uuid.Nil, false
	}

	return id, true
}

// QueryUUID extracts and parses a UUID query parameter
// Returns an error and writes an error response if parsing fails
func QueryUUID(c *gin.Context, paramName string) (uuid.UUID, bool) {
	paramStr := c.Query(paramName)
	if paramStr == "" {
		err := apperrors.NewValidationError("missing " + paramName + " query parameter")
		HandleError(c, err)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(paramStr)
	if err != nil {
		err := apperrors.NewValidationError("invalid " + paramName + ": must be a valid UUID")
		HandleError(c, err)
		return uuid.Nil, false
	}

	return id, true
}
