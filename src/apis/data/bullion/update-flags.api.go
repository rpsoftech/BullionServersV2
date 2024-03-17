package bullion

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type apiUpdateBullionFlagsBody struct {
	BullionId   string                     `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
	FlagsEntity *interfaces.FlagsInterface `bson:"flagsEntity" json:"flagsEntity" validate:"required"`
}

func apiUpdateBullionFlags(c *fiber.Ctx) error {
	body := new(apiUpdateBullionFlagsBody)
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
	if body.FlagsEntity.BullionId != body.BullionId {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Bullion Id Mismatch",
			Name:       "INVALID_INPUT",
		}
	}
	entity, err := services.FlagService.UpdateFlags(body.FlagsEntity, userID)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}
