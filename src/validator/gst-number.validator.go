package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var gstRegex, _ = regexp.Compile(`\d{2}[A-Z]{5}\d{4}[A-Z]{1}[A-Z\d]{1}[Z]{1}[A-Z\d]{1}`)

func validateGstNumber(fl validator.FieldLevel) bool {
	gstnumber := fl.Field().String()
	return gstRegex.MatchString(gstnumber)
}
