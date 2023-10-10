package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetBullionDetailsByShortName(c *fiber.Ctx) error {
	var body map[string]string = c.Queries()
	shortName, ok := body["name"]
	if !ok || len(shortName) < 2 {
		err := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Invalid Short Name",
			Name:       "ERROR_INVALID_INPUT",
		}
		return err
	}
	if entity, err := services.BullionDetailsService.GetBullionDetailsByShortName(shortName); err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}

func apiGetBullionDetailsById(c *fiber.Ctx) error {
	var body map[string]string = c.Queries()
	bullionId, ok := body["id"]
	if !ok || len(bullionId) < 2 {
		err := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Invalid Short Name",
			Name:       "ERROR_INVALID_INPUT",
		}
		return err
	}
	if entity, err := services.BullionDetailsService.GetBullionDetailsByBullionId(bullionId); err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
