package middleware

import (
	"strconv"

	apperrors "restaurant/internal/errors"

	"github.com/gin-gonic/gin"
)

// ValidateQueryPagination validates offset and limit query parameters
// Returns validated offset and limit, or error if invalid
func ValidateQueryPagination(c *gin.Context, defaultLimit int) (int, int, bool) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		err := apperrors.NewValidationError("offset must be a non-negative integer")
		HandleError(c, err)
		return 0, 0, false
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		err := apperrors.NewValidationError("limit must be between 1 and 100")
		HandleError(c, err)
		return 0, 0, false
	}

	return offset, limit, true
}

// GetIntQueryParam safely extracts and validates an integer query parameter
func GetIntQueryParam(c *gin.Context, paramName string, required bool) (int, bool) {
	value := c.Query(paramName)
	if value == "" {
		if required {
			err := apperrors.NewValidationError(paramName + " query parameter is required")
			HandleError(c, err)
			return 0, false
		}
		return 0, true
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		errResp := apperrors.NewValidationError(paramName + " must be a valid integer")
		HandleError(c, errResp)
		return 0, false
	}

	return intVal, true
}

// GetStringQueryParam safely extracts a string query parameter
func GetStringQueryParam(c *gin.Context, paramName string, required bool) (string, bool) {
	value := c.Query(paramName)
	if value == "" && required {
		err := apperrors.NewValidationError(paramName + " query parameter is required")
		HandleError(c, err)
		return "", false
	}
	return value, true
}
