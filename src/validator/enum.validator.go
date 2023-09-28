package validator

import (
	"github.com/go-playground/validator/v10"
)

type enumValidationFunction func(value string) bool

var validatorMap = map[string]enumValidationFunction{}

func validateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	methodName := fl.Param()
	if methodName == "" {
		panic("Method name is not passed in params for enum validator")
	}

	function, found := validatorMap[methodName]
	if !found {
		panic("Please Register enum validator function for " + methodName)
	}
	ok := function(value)

	if !ok {
		panic("Value " + value + " is not valid for enum " + fl.Parent().String() + "." + fl.FieldName())
	}
	return true
}

func RegisterEnumValidatorFunc(name string, function enumValidationFunction) {
	validatorMap[name] = function
}
