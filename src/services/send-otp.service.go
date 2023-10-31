package services

import "github.com/rpsoftech/bullion-server/src/redis"

type sendOtpServiceStruct struct {
	redisRepo *redis.RedisClientStruct
	eventBus  *eventBusService
}

var SendOtpService *sendOtpServiceStruct

func InitAndGetSendOtpService() *sendOtpServiceStruct {
	if SendOtpService == nil {
		SendOtpService = &sendOtpServiceStruct{
			redisRepo: redis.InitRedisAndRedisClient(),
			eventBus:  getEventBusService(),
		}
	}
	return SendOtpService
}

// func (s *sendOtpServiceStruct)
