package mongo

import (
	"os"

	"github.com/rpsoftech/bullion-server/src/env"
)

func init() {
	DbUrl := os.Getenv(env.DB_URL_KEY)
	if DbUrl == "" {
		panic("Pleas pass valid DB Url")
	}
}
