package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiAddNewProductBody struct {
	*interfaces.ProductBaseStruct  `json:"product" validate:"required"`
	*interfaces.CalcSnapshotStruct `json:"calcSnapShot" validate:"required"`
}

func apiAddNewProduct(c *fiber.Ctx) error {
	body := new(apiAddNewProductBody)
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
	entity, err := services.ProductService.AddNewProduct(body.ProductBaseStruct, body.CalcSnapshotStruct, userId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
