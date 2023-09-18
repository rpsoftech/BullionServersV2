package validator

import (
	v "github.com/go-playground/validator/v10"
)

var Validator = v.New()

func init() {
	println("Registerd")
	Validator.RegisterValidation("port", ValidatePort)
	// Validator.Re("port", ValidatePort)
}
