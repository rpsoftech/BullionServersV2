package bullion

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddBullionApis(router fiber.Router) {
	// app.Get("/bullions", getBullions)
	{
		adminGroup := router.Group("/admin", middleware.AllowOnlyBigAdmins.Validate)
		adminGroup.Patch("/updateBullionFlags", apiUpdateBullionFlags)
	}
	{
		tradeUserGroup := router.Group("/user")
		tradeUserGroup.Get("/getBullionFlags", apiGetBullionFlags)
	}
}
