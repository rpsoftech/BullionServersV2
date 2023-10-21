package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiGetFeedsApiPaginatedBody struct {
	BullionId string `json:"bullionId" validate:"required,uuid"`
	Page      int64  `json:"page" validate:"required,min=1"`
	Limit     int64  `json:"limit" validate:"required,min=1"`
}

func apiGetFeedsApiPaginated(c *fiber.Ctx) error {
	body := new(apiGetFeedsApiPaginatedBody)
	c.QueryParser(body)
	if err := interfaces.ValidateBullionIdMatchingInToken(c, body.BullionId); err != nil {
		return err
	}
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	entity, err := services.FeedsService.FetchPaginatedFeedsByBullionId(body.BullionId, body.Page-1, body.Limit)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
