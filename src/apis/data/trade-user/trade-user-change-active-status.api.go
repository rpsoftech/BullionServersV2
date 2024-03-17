package tradeuser

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiChangeTradeUserStatusBody struct {
	ID        string `bson:"id" json:"id" validate:"required,uuid"`
	IsActive  bool   `bson:"isActive" json:"isActive" validate:"boolean"`
	BullionId string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
}

func apiChangeTradeUserStatus(c *fiber.Ctx) error {
	body := new(apiChangeTradeUserStatusBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	err = services.TradeUserService.TradeUserChangeStatus(body.ID, body.BullionId, body.IsActive, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(&fiber.Map{
			"success": true,
			"status":  body.IsActive,
		})
	}

}
