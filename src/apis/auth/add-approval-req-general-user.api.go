package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type addApprovalReqGeneralUserBody struct {
	getGeneralUserBody
	BullionId string `json:"bullionId" validate:"required"`
}

func apiSendApprovalReqGeneralUser(c *fiber.Ctx) error {
	body := new(addApprovalReqGeneralUserBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.GeneralUserService.CreateApprovalRequest(body.Id, body.Password, body.BullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
