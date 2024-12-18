package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetLiveRate(c *fiber.Ctx) error {
	return c.JSON(services.LiveRateService.LastRateMap)
}
