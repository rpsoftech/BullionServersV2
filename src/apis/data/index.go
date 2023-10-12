package data

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddDataPackage(router fiber.Router) {
	router.Use(middleware.AllowOnlyValidTokenMiddleWare)
	router.Use(middleware.AllowAllUsers.Validate)
	router.Get("AAAAa", func(c *fiber.Ctx) error {
		return c.SendString("kjAIsoaSHIUSh")
	})
	// router.Get("/products", data.GetProducts)
	// router.Get("/deviceId", generateDeviceId)
	// router.Get("/bullion-details-by-short-name", apiGetBullionDetailsByShortName)
	// router.Get("/bullion-details-by-id", apiGetBullionDetailsById)
	// {
	// 	generalUserGroup := router.Group("general-user")
	// 	generalUserGroup.Post("/register", apiRegisterNewGeneralUser)
	// 	generalUserGroup.Get("/get", apiGetGeneralUserDetailsByIdPassword)
	// 	generalUserGroup.Post("/send-for-approval", apiSendApprovalReqGeneralUser)
	// 	generalUserGroup.Post("/get-general-user-token", apiGetGeneralUserToken)
	// 	generalUserGroup.Post("/refresh-token", apiGeneralUSerRefreshToken)
	// }
	// {
	// 	adminAuthGroup := router.Group("admin")
	// 	adminAuthGroup.Post("/login", apiAdminLogin)
	// }
}
