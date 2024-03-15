package tradeusergroup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiUpdateTradeUserGroupMapBody struct {
	Base      *[]interfaces.TradeUserGroupMapEntity `bson:"base" json:"base" validate:"required"`
	GroupId   string                                `bson:"groupId" json:"groupId" validate:"required,uuid"`
	BullionId string                                `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
}

func apiUpdateTradeUserGroupMap(c *fiber.Ctx) error {
	body := new(apiUpdateTradeUserGroupMapBody)
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
	entity, err := services.TradeUserGroupService.UpdateTradeGroupMap(body.Base, body.GroupId, body.BullionId, userID)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}

}
