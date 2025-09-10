package lib

import (
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
}

// ValidateStruct validates a struct using the global validator
func ValidateStruct(s interface{}) error {
	return Validator.Struct(s)
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return Validator.Var(field, tag)
}
