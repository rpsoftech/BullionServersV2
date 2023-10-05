package auth

import (
	"github.com/gofiber/fiber/v2"
)

func AddAuthPackages(router fiber.Router) {
	// router.Use(middleware.AllowAllUsers.Validate)
	router.Get("/deviceId", generateDeviceId)
	{
		generalUserGroup := router.Group("general-user")
		generalUserGroup.Post("/register", apiRegisterNewGeneralUser)
		generalUserGroup.Get("/get", apiGetGeneralUserDetailsByIdPassword)
		generalUserGroup.Post("/send-for-approval", apiSendApprovalReqGeneralUser)
		generalUserGroup.Post("/get-general-user-token", apiGetGeneralUserToken)
		generalUserGroup.Post("/refresh-token", apiGeneralUSerRefreshToken)
	}
}
