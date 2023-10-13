package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiAddNewProductBody struct {
	interfaces.ProductBaseStruct  `json:"product" validate:"required"`
	interfaces.CalcSnapshotStruct `json:"calcSnapShot" validate:"required"`
}

func apiAddNewProduct(c *fiber.Ctx) error {
	body := new(apiAddNewProductBody)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	return c.JSON(body)
}
