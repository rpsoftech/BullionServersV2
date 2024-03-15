package tradeusergroup

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetTradeGroupDetailsByBullionId(c *fiber.Ctx) error {
	id := c.Query("bullionId")
	if id == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Bullion Id",
			Name:       "INVALID_INPUT",
		}
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, id); err != nil {
		return err
	}
	entity, err := services.TradeUserGroupService.GetAllGroupsByBullionId(id)
	if err != nil {
		return err
	} else {
		return c.JSON(fiber.Map{
			"total": len(*entity),
			"data":  entity,
		})
	}
}

func apiGetTradeGroupDetailsById(c *fiber.Ctx) error {
	id := c.Query("groupId")
	if id == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Group ID",
			Name:       "PLEASE_PASS_VALID_GROUP_ID",
		}
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserGroupService.GetGroupByGroupId(id, bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}

func apiGetTradeGroupMapDetailsByGroupId(c *fiber.Ctx) error {
	id := c.Query("groupId")
	if id == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Group ID",
			Name:       "PLEASE_PASS_VALID_GROUP_ID",
		}
	}
	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	entity, err := services.TradeUserGroupService.GetGroupMapByGroupId(id, bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(fiber.Map{
			"total": len(*entity),
			"data":  entity,
		})
	}
}
