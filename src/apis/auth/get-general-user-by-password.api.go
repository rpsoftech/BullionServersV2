package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type getGeneralUserBody struct {
	Id       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

func apiGetGeneralUserDetailsByIdPassword(c *fiber.Ctx) error {
	body := new(getGeneralUserBody)
	c.QueryParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.GeneralUserService.GetGeneralUserDetailsByIdPassword(body.Id, body.Password)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
