package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Init initializes the validator
func Init() {
	once.Do(func() {
		validate = validator.New()
	})
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}