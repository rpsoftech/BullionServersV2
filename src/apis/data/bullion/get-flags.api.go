package bullion

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetBullionFlags(c *fiber.Ctx) error {
	bullionId := c.Query("bullionId")
	if bullionId == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "bullionId is required",
			Name:       "INVALID_INPUT",
		}
	}

	if err := interfaces.ValidateBullionIdMatchingInToken(c, bullionId); err != nil {
		return err
	}

	entity, err := services.FlagService.GetFlags(bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
