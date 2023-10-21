package auth

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGeneralUSerRefreshToken(c *fiber.Ctx) error {
	var body map[string]string
	json.Unmarshal(c.Body(), &body)
	token, ok := body["token"]
	if !ok {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Token",
			Name:       "INVALID_INPUT",
		}
	}
	entity, err := services.GeneralUserService.RefreshToken(token)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
