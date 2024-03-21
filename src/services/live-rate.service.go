package services

import (
	"encoding/json"
	"time"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/redis"
)

type LiveRateService struct {
	redisRepo   *redis.RedisClientStruct
	LastRateMap interfaces.LiveRateData
}

var LiveRateServiceInstance *LiveRateService

func init() {
	service := getLiveRateService()
	service.lastRateReaderFromRedis()
}

func getLiveRateService() *LiveRateService {
	if LiveRateServiceInstance == nil {
		LiveRateServiceInstance = &LiveRateService{
			redisRepo:   redis.InitRedisAndRedisClient(),
			LastRateMap: make(map[interfaces.SymbolsEnum]map[interfaces.PriceKeyEnum]float64),
		}
		println("Live Rate Service Initialized")
	}
	return LiveRateServiceInstance
}

func (s *LiveRateService) GetLastRate() *interfaces.LiveRateData {
	return &s.LastRateMap
}

func (s *LiveRateService) lastRateReaderFromRedis() {
	println("Reading Last Rate From Redis Started")
	go func() {
		for {
			data, err := s.redisRepo.GetByteData("LastRate")
			// s.LastRateMap = res
			if err != nil && len(data) > 0 {
				json.Unmarshal(data, &s.LastRateMap)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
