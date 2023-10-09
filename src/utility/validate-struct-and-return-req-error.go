package utility

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/validator"
)

func ValidateReqInput(body interface{}) (err *interfaces.RequestError) {
	err = &interfaces.RequestError{
		StatusCode: 400,
		Code:       interfaces.ERROR_INVALID_INPUT,
		Message:    "",
		Name:       "INVALID_INPUT",
		Extra:      nil,
	}
	return ValidateStructAndReturnReqError(body, err)
}

func ValidateStructAndReturnReqError(data interface{}, err *interfaces.RequestError) *interfaces.RequestError {
	if errs := validator.Validator.Validate(data); len(errs) > 0 {
		err.Extra = errs
		err.AppendValidationErrors(errs)
		return err
	} else {
		return nil
	}
}
