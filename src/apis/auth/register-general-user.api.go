package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
)

type registerGeneralUserBody struct {
	BullionId string      `json:"bullionId" validate:"required"`
	User      interface{} `json:"user" validation:"required"`
}

func apiRegisterNewGeneralUser(c *fiber.Ctx) (err error) {
	body := new(registerGeneralUserBody)
	c.BodyParser(body)
	services.GeneralUserService.RegisterNew(body.BullionId)
	return c.JSON(body)
}
