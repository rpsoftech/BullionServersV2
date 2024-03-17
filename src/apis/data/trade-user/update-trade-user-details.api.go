package tradeuser

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiUpdateTradeUserDetails(c *fiber.Ctx) error {
	body := new(interfaces.TradeUserEntity)
	c.BodyParser(body)
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	err = services.TradeUserService.UpdateTradeUser(body, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(&fiber.Map{
			"success": true,
		})
	}

}
