package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddAuthPackages(router fiber.Router) {
	// router.Use(middleware.AllowAllUsers.Validate)
	router.Get("/deviceId", generateDeviceId)
	router.Get("/bullion-details-by-short-name", apiGetBullionDetailsByShortName)
	router.Get("/bullion-details-by-id", apiGetBullionDetailsById)
	{
		generalUserGroup := router.Group("general-user")
		generalUserGroup.Post("/register", apiRegisterNewGeneralUser)
		generalUserGroup.Get("/get", apiGetGeneralUserDetailsByIdPassword)
		generalUserGroup.Post("/send-for-approval", apiSendApprovalReqGeneralUser)
		generalUserGroup.Post("/get-general-user-token", apiGetGeneralUserToken)
		generalUserGroup.Post("/refresh-token", apiGeneralUSerRefreshToken)
	}
	{
		adminAuthGroup := router.Group("admin")
		adminAuthGroup.Post("/login", apiAdminLogin)
		// adminApiGroup := adminAuthGroup.Use(middleware.AllowOnlyBigAdmins.Validate)

	}
	{
		tradeUserGroup := router.Group("trade-user")
		tradeUserGroup.Use(middleware.AllowOnlyValidTokenMiddleWare)
		tradeUserGroup.Use(middleware.AllowAllUsers.Validate)
		tradeUserGroup.Post("register", apiTradeUserRegister)
		tradeUserGroup.Post("resend-otp", apiTradeUserResendOtp)
		tradeUserGroup.Put("verify-otp", apiTradeUserVerifyOtp)
	}
}
