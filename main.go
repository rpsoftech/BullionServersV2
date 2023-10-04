package main

import (
	"strconv"
	"time"

	j "github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/apis"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/middleware"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility/firebase"
	"github.com/rpsoftech/bullion-server/src/utility/jwt"
)

func deferMainFunc() {
	println("Closing...")
	mongodb.DeferFunction()
}

func main() {
	defer deferMainFunc()
	firebase.Init()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			mappedError, ok := err.(*interfaces.RequestError)
			if !ok {
				println(err.Error())
				return c.Status(500).JSON(interfaces.RequestError{
					Code:    interfaces.ERROR_INTERNAL_ERROR,
					Message: "Some Internal Error",
					Name:    "Global Error Handler Function",
				})
			}
			return c.Status(mappedError.StatusCode).JSON(mappedError)
		},
	})
	app.Use(middleware.TokenDecrypter)
	app.Get("/token", func(c *fiber.Ctx) error {
		a, _ := services.AccessTokenService.GenerateToken(jwt.CustomClaims{
			Claims: interfaces.UserRolesInterface{
				Role: interfaces.ROLE_ADMIN,
			},
			RegisteredClaims: j.RegisteredClaims{
				IssuedAt: j.NewNumericDate(time.Now()),
			},
		})
		return c.SendString(a)
	})
	// repos.BullionSiteInfoRepo.Save(interfaces.CreateNewBullionSiteInfo("Akshat Bullion", "https://akshatbullion.com").AddGeneralUserInfo(true, true))
	// app.Get("/", func(c *fiber.Ctx) error {
	// bull := repos.BullionSiteInfoRepo.FindOne("ad3cee16-e8d7-4a27-a060-46d99c133273")
	// return c.JSON(bull)

	// return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatgold.com"))
	// return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatbullion.com"))
	// return c.SendString("Hello, World!")
	// })
	apis.AddApis(app.Group("/v1"))
	hostAndPort := ""
	if env.Env.APP_ENV == env.APP_ENV_LOCAL || env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		hostAndPort = "127.0.0.1"
	}
	hostAndPort = hostAndPort + ":" + strconv.Itoa(env.Env.PORT)
	app.Listen(hostAndPort)
	// log.Fatal(app.Listen(":" + strconv.Itoa(env.Env.PORT)))
}

/*

func main() {
	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	// Custom struct validation tag format
	myValidator.validator.RegisterValidation("teener", func(fl validator.FieldLevel) bool {
		// User.Age needs to fit our needs, 12-18 years old.
		return fl.Field().Int() >= 12 && fl.Field().Int() <= 18
	})

	app.Get("/", func(c *fiber.Ctx) error {
		user := &User{
			Name: c.Query("name"),
			Age:  c.QueryInt("age"),
		}

		// Validation
		if errs := myValidator.Validate(user); len(errs) > 0 && errs[0].Error {
			errMsgs := make([]string, 0)

			for _, err := range errs {
				errMsgs = append(errMsgs, fmt.Sprintf(
					"[%s]: '%v' | Needs to implement '%s'",
					err.FailedField,
					err.Value,
					err.Tag,
				))
			}

			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: strings.Join(errMsgs, " and "),
			}
		}

		// Logic, validated with success
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
*/
