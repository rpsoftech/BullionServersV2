package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddProduct(router fiber.Router) {
	router.Get("/getAll", apiGetProducts)
	router.Get("/getProduct", apiGetProducts)
	adminGroup := router.Group("", middleware.AllowAllAdmins.Validate)
	adminGroup.Post("/add", middleware.AllowAllAdmins.Validate, apiAddNewProduct)
	adminGroup.Post("/update", middleware.AllowAllAdmins.Validate, apiUpdateProducts)
	adminGroup.Post("/updateCalcSnapShot", middleware.AllowAllAdmins.Validate, apiUpdateProductCalcSnapshot)
	adminGroup.Post("/updateSequence", middleware.AllowAllAdmins.Validate, apiUpdateProductSequence)
}
