package data

import (
	"github.com/gofiber/fiber/v2"
	bankdetails "github.com/rpsoftech/bullion-server/src/apis/data/bank-details"
	"github.com/rpsoftech/bullion-server/src/apis/data/feeds"
	"github.com/rpsoftech/bullion-server/src/apis/data/product"
	tradeuser "github.com/rpsoftech/bullion-server/src/apis/data/trade-user"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddDataPackage(router fiber.Router) {
	router.Use(middleware.AllowOnlyValidTokenMiddleWare)
	router.Use(middleware.AllowAllUsers.Validate)
	{
		productGroup := router.Group("/product")
		product.AddProduct(productGroup)
	}
	{
		feedGroup := router.Group("/feeds")
		feeds.AddFeedsAndNotificationSection(feedGroup)
	}
	{
		bankGroup := router.Group("/bank-details")
		bankdetails.AddBankDetailsAPIs(bankGroup)
	}
	{
		tradeUserGroup := router.Group("/trade-user")
		tradeuser.AddTradeUserAPIs(tradeUserGroup)
	}
}
