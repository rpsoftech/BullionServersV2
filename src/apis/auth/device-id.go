package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateDeviceId(c *fiber.Ctx) error {
	id := uuid.New()
	return c.SendString(strings.ReplaceAll(id.String(), "-", ""))
	// return c.JSON(fiber.Map{
	// 	"id": strings.ReplaceAll(id.String(), "-", ""),
	// })
}
