package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

func apiUpdateProducts(c *fiber.Ctx) error {
	body := new([]interfaces.UpdateProductApiBody)
	c.BodyParser(body)
	if len(*body) == 0 {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please pass at least 1 product to update",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}

	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	for _, ele := range *body {
		if err := utility.ValidateReqInput(&ele); err != nil {
			return err
		}
		if ele.BullionId != bullionId {
			return &interfaces.RequestError{
				StatusCode: 403,
				Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
				Message:    "You can not change other resources",
				Name:       "ERROR_BULLION_ID_MISMATCH",
				Extra: fiber.Map{
					"bullionId": bullionId,
					"product":   ele,
				},
			}
		}
	}
	if entity, err := services.ProductService.UpdateProduct(body, userId); err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
func apiUpdateProductCalcSnapshot(c *fiber.Ctx) error {
	body := new([]interfaces.UpdateProductCalcSnapshotApiBody)
	c.BodyParser(body)
	if len(*body) == 0 {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please pass at least 1 product to update",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}

	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	for _, ele := range *body {
		if err := utility.ValidateReqInput(&ele); err != nil {
			return err
		}
		if ele.BullionId != bullionId {
			return &interfaces.RequestError{
				StatusCode: 403,
				Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
				Message:    "You can not change other resources",
				Name:       "ERROR_BULLION_ID_MISMATCH",
				Extra: fiber.Map{
					"bullionId": bullionId,
					"product":   ele,
				},
			}
		}
	}
	if entity, err := services.ProductService.UpdateProductCalcSnapshot(body, userId); err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}

func apiUpdateProductSequence(c *fiber.Ctx) error {
	body := new([]interfaces.UpdateProductCalcSequenceApiBody)
	c.BodyParser(body)
	if len(*body) < 2 {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please pass at least 2 product to update",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	userId, err := interfaces.ExtractTokenUserIdFromCtx(c)
	if err != nil {
		return err
	}

	bullionId, err := interfaces.ExtractBullionIdFromCtx(c)
	if err != nil {
		return err
	}
	for _, ele := range *body {
		if err := utility.ValidateReqInput(&ele); err != nil {
			return err
		}
		if ele.BullionId != bullionId {
			return &interfaces.RequestError{
				StatusCode: 403,
				Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
				Message:    "You can not change other resources",
				Name:       "ERROR_BULLION_ID_MISMATCH",
				Extra: fiber.Map{
					"bullionId": bullionId,
					"product":   ele,
				},
			}
		}
	}
	if entity, err := services.ProductService.UpdateProductSequence(body, userId, bullionId); err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
