package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiSendFeedAsNotification(c *fiber.Ctx) error {
	body := new(interfaces.FeedsBase)
	c.BodyParser(body)
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}
	err = services.FeedsService.SendNotification(body.BullionId, body, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(utility.SuccessResponse)
	}
}
