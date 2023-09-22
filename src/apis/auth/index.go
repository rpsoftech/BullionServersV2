package auth

import "github.com/gofiber/fiber/v2"

func AddAuthPackages(router fiber.Router) {
	router.Get("/deviceId", generateDeviceId)
}
