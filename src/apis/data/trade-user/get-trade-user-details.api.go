package tradeuser

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiGetTradeUserDetailsBody struct {
	ID        string `bson:"id" json:"id" validate:"required,uuid"`
	BullionId string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
}

func apiGetTradeUserDetails(c *fiber.Ctx) error {
	body := new(apiGetTradeUserDetailsBody)
	c.QueryParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	entity, err := services.TradeUserService.FindOneUserById(body.ID)
	if entity.BullionId != body.BullionId {
		return &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "Bullion Id Mismatch For Trade User",
			Name:       "BULLION_ID_MISMATCH_FOR_TRADE_USER",
		}
	}
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}

}
