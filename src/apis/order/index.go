package order

import (
	"github.com/gofiber/fiber/v2"
	adminorder "github.com/rpsoftech/bullion-server/src/apis/order/admin-order"
	"github.com/rpsoftech/bullion-server/src/apis/order/user"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddOrderPackage(router fiber.Router) {
	router.Use(middleware.AllowOnlyValidTokenMiddleWare)
	adminorder.AddAdminOrderRouter(router.Group("/admin", middleware.AllowOnlyBigAdmins.Validate))
	user.AddUserOrderApis(router.Group("/user", middleware.AllowAllAdminsAndTradeUsers.Validate))
}
