package tradeusergroup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiUpdateTradeUserGroupBody struct {
	Base    *interfaces.TradeUserGroupBase `bson:"base" json:"base" validate:"required"`
	GroupId string                         `bson:"groupId" json:"groupId" validate:"required,uuid"`
}

func apiUpdateTradeUserGroup(c *fiber.Ctx) error {
	body := new(apiUpdateTradeUserGroupBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.Base.BullionId); err != nil {
		return err
	}
	userID, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserGroupService.UpdateTradeGroup(body.Base, body.GroupId, userID)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}

}
