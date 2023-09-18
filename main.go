package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rpsoftech/bullion-server/src/env"
)

func init() {
	godotenv.Load()
}

func main() {
	port := os.Getenv(env.PORTKey)
	if port == "" {
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
