package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a standardized application error with HTTP status code
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// StatusCode returns the HTTP status code for the error
func (e *AppError) StatusCode() int {
	return e.Code
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Standard error variables with HTTP status codes
var (
	// 400 Bad Request
	ErrInvalidRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid request",
	}

	ErrInvalidInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid input provided",
	}

	ErrValidationFailed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "validation failed",
	}

	// 404 Not Found
	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "resource not found",
	}

	ErrMenuItemNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "menu item not found",
	}

	ErrOrderNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "order not found",
	}

	ErrSessionNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "session not found",
	}

	ErrCategoryNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "category not found",
	}

	ErrOrderItemNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "order item not found",
	}

	// 409 Conflict
	ErrConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "resource already exists",
	}

	ErrDuplicateCategory = &AppError{
		Code:    http.StatusConflict,
		Message: "category already exists",
	}

	ErrDuplicateMenuItem = &AppError{
		Code:    http.StatusConflict,
		Message: "menu item already exists",
	}

	// 400 Bad Request - Business Logic
	ErrOutOfStock = &AppError{
		Code:    http.StatusBadRequest,
		Message: "item is out of stock",
	}

	ErrInsufficientStock = &AppError{
		Code:    http.StatusBadRequest,
		Message: "insufficient stock available",
	}

	ErrInvalidTransition = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid state transition",
	}

	ErrInvalidOrderStatus = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid order status",
	}

	ErrInvalidSessionStatus = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid session status",
	}

	// 500 Internal Server Error
	ErrInternal = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	}

	ErrDatabaseError = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "database operation failed",
	}

	// Foreign key constraint
	ErrForeignKeyViolation = &AppError{
		Code:    http.StatusConflict,
		Message: "cannot delete resource: it is referenced by other records",
	}
)

// NewAppError creates a new AppError with custom message
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a 404 error with custom message
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// NewValidationError creates a 400 validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// NewConflictError creates a 409 conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// NewInternalError creates a 500 internal error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError or returns nil
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// WrapError wraps an error with context and AppError type
func WrapError(code int, message string, err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		// If already an AppError, enhance with context
		return &AppError{
			Code:    appErr.Code,
			Message: fmt.Sprintf("%s: %s", message, appErr.Message),
			Err:     appErr.Err,
		}
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// GetHTTPStatusCode returns the HTTP status code for an error
func GetHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}

	// Default to 500 for unknown errors
	return http.StatusInternalServerError
}

// GetErrorMessage returns a user-friendly error message
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Message
	}

	return "internal server error"
}
