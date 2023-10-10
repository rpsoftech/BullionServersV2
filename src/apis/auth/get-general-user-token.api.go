package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiGetGeneralUserToken(c *fiber.Ctx) error {
	body := new(addApprovalReqGeneralUserBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.GeneralUserService.ValidateApprovalAndGenerateToken(body.Id, body.Password, body.BullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
