package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

func deferMainFunc() {
	println("Closing...")
	mongodb.DeferFunction()
}

func main() {
	defer deferMainFunc()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(repos.BullionSiteInfoRepo.FindOne("14a8c0fd-1d17-4421-af5e-8e12a5115361"))
		// return c.SendString("Hello, World!")
	})
	app.Listen(":" + strconv.Itoa(env.Env.PORT))

}
