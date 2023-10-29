package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rpsoftech/bullion-server/src/env"
)

var RedisClient *redis.Client
var RedisCTX = context.Background()

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     env.Env.REDIS_DB_URL,
		Password: env.Env.REDIS_DB_PASSWORD, // no password set
		DB:       env.Env.REDIS_DB_DATABASE, // use default DB
	})

	RedisClient = client
	go func() {
		err := RedisClient.Ping(RedisCTX)
		if err != nil {
			panic(err)
		}
	}()
}

func DeferFunction() {
	if err := RedisClient.Close(); err != nil {
		panic(err)
	}
}
