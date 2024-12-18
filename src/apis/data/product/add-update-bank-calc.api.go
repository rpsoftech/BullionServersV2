package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiAddUpdateBankCalcBody struct {
	GOLD_SPOT   *interfaces.BankRateCalcBase `bson:"goldSpot" json:"goldSpot" validate:"required"`
	SILVER_SPOT *interfaces.BankRateCalcBase `bson:"silverSpot" json:"silverSpot" validate:"required"`
}

func apiAddUpdateBankCalc(c *fiber.Ctx) error {
	body := new(apiAddUpdateBankCalcBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.BankRateCalcService.SaveBankRateCalc(body.GOLD_SPOT, body.SILVER_SPOT, bullionId, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
