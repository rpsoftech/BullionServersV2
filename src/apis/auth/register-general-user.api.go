package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/validator"
)

type registerGeneralUserBody struct {
	BullionId string      `json:"bullionId" validate:"required"`
	User      interface{} `json:"user" validate:"required"`
}

func apiRegisterNewGeneralUser(c *fiber.Ctx) error {
	body := new(registerGeneralUserBody)
	c.BodyParser(body)
	if errs := validator.Validator.Validate(body); len(errs) > 0 {
		err := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "",
			Name:       "INVALID_INPUT",
			Extra:      errs,
		}
		err.AppendValidationErrors(errs)
		return err
	}
	entity, err := services.GeneralUserService.RegisterNew(body.BullionId, body.User)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
