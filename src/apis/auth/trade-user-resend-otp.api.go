package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiTradeUserResendOtp(c *fiber.Ctx) error {
	// var body
	body := new(map[string]string)
	json.Unmarshal(c.Body(), body)
	token, ok := (*body)["token"]
	if !ok {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Token",
			Name:       "INVALID_INPUT",
		}
	}

	err := services.TradeUserService.VerifyTokenAndResendOTP(token)
	if err != nil {
		return err
	} else {
		return c.JSON(utility.SuccessResponse)
	}
}
