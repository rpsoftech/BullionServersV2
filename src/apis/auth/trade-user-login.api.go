package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiTradeUserLoginNumberBody struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

func apiTradeUserLoginNumber(c *fiber.Ctx) error {
	// var body
	body := new(apiTradeUserLoginNumberBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserService.LoginWithNumberAndPassword(body.Number, body.Password, bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
