package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddFeedsAndNotificationSection(router fiber.Router) {
	router.Get("/getPaginated", apiGetFeedsApiPaginated)
	{
		adminGroup := router.Group("", middleware.AllowOnlyBigAdmins.Validate)
		adminGroup.Post("/sendNotification", apiSendFeedAsNotification)
		adminGroup.Post("/add", apiAddNewFeed)
		adminGroup.Post("/update", apiUpdateFeed)
		adminGroup.Delete("/delete", apiDeleteFeeds)
	}
}
