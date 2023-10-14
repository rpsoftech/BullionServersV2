package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetProducts(c *fiber.Ctx) (err error) {

	id := c.Query("bullionId")
	if id == "" {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass Valid Bullion Id",
			Name:       "INVALID_INPUT",
		}
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, id); err != nil {
		return err
	}
	productId := c.Query("productId")
	var entity interface{}
	if productId != "" {
		entity, err = services.ProductService.GetProductsById(id, productId)
	} else {
		entity, err = services.ProductService.GetProductsByBullionId(id)
	}
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
