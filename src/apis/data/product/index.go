package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddProduct(router fiber.Router) {
	router.Get("/getAll", apiGetProducts)
	router.Get("/getProduct", apiGetProducts)
	adminGroup := router.Use(middleware.AllowAllAdmins.Validate)
	adminGroup.Put("/add", apiAddNewProduct)
	adminGroup.Patch("/update", apiUpdateProducts)
	adminGroup.Patch("/updateCalcSnapShot", apiUpdateProductCalcSnapshot)
	adminGroup.Patch("/updateSequence", apiUpdateProductSequence)
}
