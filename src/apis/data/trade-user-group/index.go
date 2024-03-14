package tradeusergroup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddTradeUserAPIs(router fiber.Router) {

	adminAndTradeGroup := router.Use(middleware.AllowAllAdminsAndTradeUsers.Validate)
	adminAndTradeGroup.Get("/getTradeGroupDetailsByID", apiGetTradeGroupDetailsById)
	adminAndTradeGroup.Get("/getTradeGroupDetailsByBullionId", apiGetTradeGroupDetailsByBullionId)
	adminAndTradeGroup.Get("/getTradeGroupMapDetailsByGroupId", apiGetTradeGroupMapDetailsByGroupId)
	adminGroup := router.Use(middleware.AllowAllAdmins.Validate)
	adminGroup.Post("/createNewTradeGroup", apiCreateNewTradeGroup)
}
