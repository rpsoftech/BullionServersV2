package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiAddNewFeed(c *fiber.Ctx) error {
	body := new(interfaces.FeedsBase)
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
	entity := &interfaces.FeedsEntity{
		BaseEntity: &interfaces.BaseEntity{},
		FeedsBase:  body,
	}
	entity.CreateNewId()
	_, err = services.FeedsService.AddAndUpdateNewFeeds(entity, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
func apiUpdateFeed(c *fiber.Ctx) error {
	body := new(interfaces.FeedUpdateRequestBody)
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
	entity, err := services.FeedsService.UpdateFeeds(body.FeedsBase, body.FeedId, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
