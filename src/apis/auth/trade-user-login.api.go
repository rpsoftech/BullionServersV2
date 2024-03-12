package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiTradeUserLoginUNumberBody struct {
	UNumber  string `bson:"uNumber" json:"uNumber" validate:"required"`
	Password string `bson:"rawPassword" json:"password" mapstructure:"password" validate:"required,min=4"`
}
type apiTradeUserLoginNumberBody struct {
	Number   string `bson:"number" json:"number" validate:"required,min=12,max=12"`
	Password string `bson:"rawPassword" json:"password" mapstructure:"password" validate:"required,min=4"`
}
type apiTradeUserLoginEmailBody struct {
	Email    string `bson:"email" json:"email" validate:"required"`
	Password string `bson:"rawPassword" json:"password" mapstructure:"password" validate:"required,min=4"`
}

func apiTradeUserLoginNumber(c *fiber.Ctx) error {
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
func apiTradeUserLoginUNumber(c *fiber.Ctx) error {
	body := new(apiTradeUserLoginUNumberBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserService.LoginWithUNumberAndPassword(body.UNumber, body.Password, bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
func apiTradeUserLoginEmail(c *fiber.Ctx) error {
	body := new(apiTradeUserLoginEmailBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserService.LoginWithEmailAndPassword(body.Email, body.Password, bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
