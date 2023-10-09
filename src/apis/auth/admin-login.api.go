package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiAdminLoginBody struct {
	UserName  string `json:"uname" validate:"required"`
	Password  string `json:"password" validate:"required"`
	BullionId string `json:"bullionId" validate:"required"`
}

func apiAdminLogin(c *fiber.Ctx) error {
	body := new(apiAdminLoginBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.AdminUserService.ValidateUserAndGenerateToken(body.UserName, body.Password, body.BullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
