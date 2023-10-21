package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddProduct(router fiber.Router) {
	router.Get("/getAll", apiGetProducts)
	router.Get("/getProduct", apiGetProducts)
	router.Post("/add", middleware.AllowOnlyAdmins.Validate, apiAddNewProduct)
	router.Post("/update", middleware.AllowOnlyAdmins.Validate, apiUpdateProducts)
	router.Post("/updateCalcSnapShot", middleware.AllowOnlyAdmins.Validate, apiUpdateProductCalcSnapshot)
	router.Post("/updateSequence", middleware.AllowOnlyAdmins.Validate, apiUpdateProductSequence)
}
