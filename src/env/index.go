package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rpsoftech/bullion-server/src/validator"
)

type EnvInterface struct {
	PORT    int    `json:"PORT" validate:"required,port"`
	DB_URL  string `json:"DB_URL" validate:"required,url"`
	DB_NAME string `json:"DB_NAME_KEY" validate:"required,min=3"`
}

var Env *EnvInterface

func init() {
	godotenv.Load()
	PORT, err := strconv.Atoi(os.Getenv(PORTKey))
	if err != nil {
		panic("Please Pass Valid Port")
	}
	Env = &EnvInterface{
		PORT:    PORT,
		DB_NAME: os.Getenv(DB_NAME_KEY),
		DB_URL:  os.Getenv(DB_URL_KEY),
	}
	err = validator.Validator.Struct(Env)
	if err != nil {
		panic(err)
	}
}
