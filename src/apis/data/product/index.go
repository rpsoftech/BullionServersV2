package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddProduct(router fiber.Router) {
	router.Get("/getAll", apiGetProducts)
	router.Get("/getProduct", apiGetProducts)
	adminGroup := router.Group("", middleware.AllowAllAdmins.Validate)
	adminGroup.Put("/add", middleware.AllowAllAdmins.Validate, apiAddNewProduct)
	adminGroup.Patch("/update", middleware.AllowAllAdmins.Validate, apiUpdateProducts)
	adminGroup.Patch("/updateCalcSnapShot", middleware.AllowAllAdmins.Validate, apiUpdateProductCalcSnapshot)
	adminGroup.Patch("/updateSequence", middleware.AllowAllAdmins.Validate, apiUpdateProductSequence)
}
