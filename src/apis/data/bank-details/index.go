package bankdetails

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddBankDetailsAPIs(router fiber.Router) {
	router.Get("/getAll", apiGetBankDetails)
	{
		adminGroup := router.Use(middleware.AllowOnlyBigAdmins.Validate)
		adminGroup.Put("/add", apiAddNewBankDetails)
		adminGroup.Patch("/update", apiUpdateBankDetails)
		adminGroup.Delete("/delete", apiDeleteBankDetails)
	}
}
