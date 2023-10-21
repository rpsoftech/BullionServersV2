package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiDeleteFeeds(c *fiber.Ctx) error {
	body := new(interfaces.FeedDeleteRequestBody)
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
	entity, err := services.FeedsService.DeleteById(body.FeedId, body.BullionId, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
