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
	// repos.BullionSiteInfoRepo.Save(interfaces.CreateNewBullionSiteInfo("Akshat Bullion", "https://akshatbullion.com").AddGeneralUserInfo(true, true))
	app.Get("/", func(c *fiber.Ctx) error {
		// bull := repos.BullionSiteInfoRepo.FindOne("ad3cee16-e8d7-4a27-a060-46d99c133273")
		// return c.JSON(bull)

		return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatgold.com"))
		// return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatbullion.com"))
		// return c.SendString("Hello, World!")
	})

	app.Listen(":" + strconv.Itoa(env.Env.PORT))
}
