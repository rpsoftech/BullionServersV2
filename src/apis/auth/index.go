package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/middleware"
)

func AddAuthPackages(router fiber.Router) {
	router.Use(middleware.AllowAllUsers.Middleware)
	router.Get("/deviceId", generateDeviceId)
}
