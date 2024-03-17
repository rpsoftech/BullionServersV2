package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiTradeUserRegister(c *fiber.Ctx) error {
	body := new(interfaces.TradeUserBase)
	c.BodyParser(body)

	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	entity, err := services.TradeUserService.VerifyAndSendOtpForNewUser(body, body.BullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(&fiber.Map{
			"token": entity,
		})
	}
}
