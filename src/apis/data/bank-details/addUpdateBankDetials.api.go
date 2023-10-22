package bankdetails

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiAddNewBankDetails(c *fiber.Ctx) error {
	body := new(interfaces.BankDetailsBase)
	c.BodyParser(body)
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.BankDetailsService.AddNewBankDetails(body, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
func apiUpdateBankDetails(c *fiber.Ctx) error {
	body := new(interfaces.UpdateBankDetailsRequestBody)
	c.BodyParser(body)
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.BankDetailsService.UpdateBankDetails(body.BankDetailsBase, body.Id, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
