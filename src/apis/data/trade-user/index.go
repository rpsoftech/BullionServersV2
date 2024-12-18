package tradeuser

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddTradeUserAPIs(router fiber.Router) {
	adminGroup := router.Use(middleware.AllowOnlyBigAdmins.Validate)
	adminGroup.Get("/getInActiveTradeUsers", apiGetInActiveTradeUsers)
	adminGroup.Get("/getTradeUser", apiGetTradeUserDetails)
	adminGroup.Patch("/updateTradeUserDetails", apiUpdateTradeUserDetails)
	adminGroup.Patch("/updateTradeUserStatus", apiChangeTradeUserStatus)
}
