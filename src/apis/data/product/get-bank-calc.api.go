package product

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetBankCalc(c *fiber.Ctx) error {

	id := c.Query("bullionId")
	if id == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Bullion Id",
			Name:       "INVALID_INPUT",
		}
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, id); err != nil {
		return err
	}
	entity, err := services.BankRateCalcService.GetBankRateCalcByBullionId(id)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
