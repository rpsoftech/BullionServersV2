package apis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/apis/auth"
	"github.com/rpsoftech/bullion-server/src/apis/data"
)

func AddApis(app fiber.Router) {
	auth.AddAuthPackages(app.Group("/auth"))
	data.AddDataPackage(app.Group("/data"))
}
