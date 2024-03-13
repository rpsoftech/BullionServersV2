package tradeusergroup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiGetTradeUserDetailsBody struct {
	Name      string `bson:"name" json:"name" validate:"required"`
	BullionId string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
}

func apiCreateNewTradeGroup(c *fiber.Ctx) error {
	body := new(apiGetTradeUserDetailsBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	userID, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserGroupService.CreateNewTradeUserGroup(body.BullionId, body.Name, userID)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}

}
