package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type registerGeneralUserBody struct {
	BullionId string      `json:"bullionId" validate:"required"`
	User      interface{} `json:"user" validate:"required"`
}

func apiRegisterNewGeneralUser(c *fiber.Ctx) error {
	body := new(registerGeneralUserBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(&body); err != nil {
		return err
	}
	entity, err := services.GeneralUserService.RegisterNew(body.BullionId, body.User)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
