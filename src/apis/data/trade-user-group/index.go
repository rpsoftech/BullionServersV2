package tradeusergroup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddTradeUserGroupAPIs(router fiber.Router) {
	// router.
	// adminGroup := router.Group("/admin", middleware.AllowOnlyBigAdmins.Validate)
	{
		adminGroup := router.Group("/admin", middleware.AllowOnlyBigAdmins.Validate)
		adminGroup.Put("/createNewTradeGroup", apiCreateNewTradeGroup)
		adminGroup.Post("/updateTradeGroup", apiUpdateTradeUserGroup)
	}
	adminAndTradeGroup := router.Use(middleware.AllowAllAdminsAndTradeUsers.Validate)
	adminAndTradeGroup.Get("/getTradeGroupDetailsByID", apiGetTradeGroupDetailsById)
	adminAndTradeGroup.Get("/getTradeGroupDetailsByBullionId", apiGetTradeGroupDetailsByBullionId)
	adminAndTradeGroup.Get("/getTradeGroupMapDetailsByGroupId", apiGetTradeGroupMapDetailsByGroupId)
}
