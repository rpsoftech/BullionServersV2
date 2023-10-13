package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddProduct(router fiber.Router) {
	router.Post("/add", middleware.AllowOnlyAdmins.Validate, apiAddNewProduct)
}
