package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rpsoftech/bullion-server/src/validator"
)

type EnvInterface struct {
	APP_ENV               AppEnv `json:"APP_ENV" validate:"required,enum=AppEnv"`
	PORT                  int    `json:"PORT" validate:"required,port"`
	DB_URL                string `json:"DB_URL" validate:"required,url"`
	DB_NAME               string `json:"DB_NAME_KEY" validate:"required,min=3"`
	ACCESS_TOKEN_KEY      string `json:"ACCESS_TOKEN_KEY" validate:"required,min=100"`
	REFRESH_TOKEN_KEY     string `json:"REFRESH_TOKEN_KEY" validate:"required,min=100"`
	FIREBASE_JSON_STRING  string `json:"FIREBASE_JSON_STRING" validate:"required"`
	FIREBASE_DATABASE_URL string `json:"FIREBASE_DATABASE_URL" validate:"required"`
}

var Env *EnvInterface

func init() {
	godotenv.Load()
	PORT, err := strconv.Atoi(os.Getenv(port_KEY))
	if err != nil {
		panic("Please Pass Valid Port")
	}
	appEnv, _ := parseAppEnv(os.Getenv(app_ENV_KEY))

	Env = &EnvInterface{
		APP_ENV:               appEnv,
		PORT:                  PORT,
		DB_NAME:               os.Getenv(db_NAME_KEY),
		DB_URL:                os.Getenv(db_URL_KEY),
		ACCESS_TOKEN_KEY:      os.Getenv(access_TOKEN_KEY),
		REFRESH_TOKEN_KEY:     os.Getenv(refresh_TOKEN_KEY),
		FIREBASE_JSON_STRING:  os.Getenv(firebase_JSON_STRING_KEY),
		FIREBASE_DATABASE_URL: os.Getenv(firebase_DATABASE_URL_KEY),
	}
	errs := validator.Validator.Validate(Env)
	if len(errs) > 0 {
		panic(err)
	}
}
