package data

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/apis/data/product"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddDataPackage(router fiber.Router) {
	router.Use(middleware.AllowOnlyValidTokenMiddleWare)
	router.Use(middleware.AllowAllUsers.Validate)
	{
		productGroup := router.Group("/product")
		product.AddProduct(productGroup)
	}
}
