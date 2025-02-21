package tradeuser

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetInActiveTradeUsers(c *fiber.Ctx) error {
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserService.FindAndReturnAllInActiveTradeUsers(bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}

}
