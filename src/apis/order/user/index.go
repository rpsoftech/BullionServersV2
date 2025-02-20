package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddUserOrderApis(router fiber.Router) {
	router.Use(middleware.AllowAllAdminsAndTradeUsers.Validate)
	// router.Put("/placeMarketOrder",)
}
