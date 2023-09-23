package validator

import (
	v "github.com/go-playground/validator/v10"
)

var Validator = v.New()

func init() {
	println("Registered")
	Validator.RegisterValidation("port", validatePort)
	Validator.RegisterValidation("gstNumber", validateGstNumber)
	Validator.RegisterValidation("enum", validateEnum)
}
