package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/events"
)

type RedisClientStruct struct {
	redisClient *redis.Client
}

var RedisClient *RedisClientStruct

var RedisCTX = context.Background()

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	// RedisClient.redisClient.Subscribe()
}

func InitRedisAndRedisClient() *RedisClientStruct {
	if RedisClient != nil {
		return RedisClient
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%d", env.Env.REDIS_DB_HOST, env.Env.REDIS_DB_PORT),
		Password: env.Env.REDIS_DB_PASSWORD, // no password set
		DB:       env.Env.REDIS_DB_DATABASE, // use default DB
	})

	RedisClient = &RedisClientStruct{
		redisClient: client,
	}
	go func() {
		res := RedisClient.redisClient.Ping(RedisCTX)
		if res.Err() != nil {
			panic(res.Err())
		}
	}()
	println("Redis Client Initialized")
	return RedisClient
}

func DeferFunction() {
	if err := RedisClient.redisClient.Close(); err != nil {
		panic(err)
	}
}

func (r *RedisClientStruct) SubscribeToChannels(channels ...string) *redis.PubSub {
	return r.redisClient.Subscribe(RedisCTX, channels...)
}

func (r *RedisClientStruct) PublishEvent(event *events.BaseEvent) {
	r.redisClient.Publish(RedisCTX, event.GetEventName(), event.GetPayloadString())
}
func (r *RedisClientStruct) GetHashValue(key string) map[string]string {
	return r.redisClient.HGetAll(RedisCTX, key).Val()
}
func (r *RedisClientStruct) GetStringData(key string) string {
	return r.redisClient.Get(RedisCTX, key).Val()
}

func (r *RedisClientStruct) RemoveKey(key ...string) {
	r.redisClient.Del(RedisCTX, key...)
}
func (r *RedisClientStruct) SetStringData(key string, value string, expiresIn int) {
	r.SetStringDataWithExpiry(key, value, time.Duration(expiresIn)*time.Second)
}
func (r *RedisClientStruct) SetStringDataWithExpiry(key string, value string, expiresIn time.Duration) {
	r.redisClient.Set(RedisCTX, key, value, expiresIn)
}
