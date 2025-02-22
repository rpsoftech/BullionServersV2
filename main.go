package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	j "github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rpsoftech/bullion-server/src/apis"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/middleware"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/redis"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility/jwt"
)

func deferMainFunc() {
	println("Closing...")
	mongodb.DeferFunction()
	redis.DeferFunction()
}

func main() {
	defer deferMainFunc()

	app := fiber.New(fiber.Config{
		ServerHeader: "Bullion Server V1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			mappedError, ok := err.(*interfaces.RequestError)
			if !ok {
				println(err.Error())
				return c.Status(500).JSON(interfaces.RequestError{
					Code:    interfaces.ERROR_INTERNAL_SERVER,
					Message: "Some Internal Error",
					Name:    "Global Error Handler Function",
				})
			}
			if mappedError.LogTheDetails {
				//Todo: Store The Details Of the Error With Body And Other Extra Details Like AUTH KEY AND ETC
			}
			return c.Status(mappedError.StatusCode).JSON(mappedError)
		},
	})
	// TODO Add middleware to recover from panics https://docs.gofiber.io/api/middleware/recover
	app.Use(logger.New())
	app.Use(middleware.TokenDecrypter)

	if env.Env.APP_ENV != env.APP_ENV_PRODUCTION {

		app.Get("/token", func(c *fiber.Ctx) error {
			a, _ := services.AccessTokenService.GenerateToken(jwt.GeneralUserAccessRefreshToken{
				Role: interfaces.ROLE_ADMIN,
				RegisteredClaims: &j.RegisteredClaims{
					IssuedAt: j.NewNumericDate(time.Now()),
				},
			})
			return c.SendString(a)
		}).Name("Temp Admin Access Token")
	}
	// repos.BullionSiteInfoRepo.Save(interfaces.CreateNewBullionSiteInfo("Akshat Bullion", "https://akshatbullion.com").AddGeneralUserInfo(true, true))
	// app.Get("/", func(c *fiber.Ctx) error {
	// bull := repos.BullionSiteInfoRepo.FindOne("ad3cee16-e8d7-4a27-a060-46d99c133273")
	// return c.JSON(bull)

	// return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatgold.com"))
	// return c.JSON(repos.BullionSiteInfoRepo.FindOneByDomain("https://akshatbullion.com"))
	// return c.SendString("Hello, World!")
	// })
	apis.AddApis(app.Group("/v1"))
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
	hostAndPort := ""
	if env.Env.APP_ENV == env.APP_ENV_LOCAL || env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		hostAndPort = "127.0.0.1"
	}
	// Print the router stack in JSON format
	data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	fmt.Println(string(data))

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
