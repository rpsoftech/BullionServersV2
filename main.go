package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/mongodb"
)

func deferMainFunc() {
	println("Closing...")
	mongodb.DeferFunction()
}

func main() {
	defer deferMainFunc()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Listen(":" + strconv.Itoa(env.Env.PORT))

}
