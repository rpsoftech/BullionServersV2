package apis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/apis/auth"
)

func AddApis(app fiber.Router) {
	auth.AddAuthPackages(app.Group("/auth"))
}
