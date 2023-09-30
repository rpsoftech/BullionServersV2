package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
)

type registerGeneralUserBody struct {
	BullionId string      `json:"bullionId" validate:"required"`
	User      interface{} `json:"user" validation:"required"`
}

func apiRegisterNewGeneralUser(c *fiber.Ctx) error {
	body := new(registerGeneralUserBody)
	c.BodyParser(body)
	entity, err := services.GeneralUserService.RegisterNew(body.BullionId, body.User)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
