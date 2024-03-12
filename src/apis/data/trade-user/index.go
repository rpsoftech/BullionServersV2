package tradeuser

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddTradeUserAPIs(router fiber.Router) {
	adminGroup := router.Use(middleware.AllowAllAdmins.Validate)
	adminGroup.Get("/getTradeUser", apiGetTradeUserDetails)
	adminGroup.Post("/updateTradeUserDetails", apiUpdateTradeUserDetails)
	adminGroup.Post("/updateTradeUserStatus", apiChangeTradeUserStatus)
}
