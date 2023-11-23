package auth

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiTradeUserVerifyOtpBody struct {
	Token string `json:"token"`
	Otp   int    `json:"otp"`
}

func apiTradeUserVerifyOtp(c *fiber.Ctx) error {
	// var body
	body := new(apiTradeUserVerifyOtpBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}

	entity, err := services.TradeUserService.VerifyTokenAndVerifyOTP(body.Token, strconv.Itoa(body.Otp))
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
