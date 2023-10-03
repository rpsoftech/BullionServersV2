package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
)

type getGeneralUserBody struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func apiGetGeneralUserDetailsByIdPassword(c *fiber.Ctx) error {
	body := new(getGeneralUserBody)
	c.QueryParser(body)
	entity, err := services.GeneralUserService.GetGeneralUserDetailsByIdPassword(body.Id, body.Password)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
