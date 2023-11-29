package services

import (
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/redis"
)

type tradeUserGroupService struct {
	redisRepo          *redis.RedisClientStruct
	eventBus           *eventBusService
	firebaseDb         *firebaseDatabaseService
	bullionService     *bullionDetailsService
	tradeUserGroupRepo *repos.TradeUserRepoStruct
}

var TradeUserGroupService *tradeUserGroupService

func getTradeUserGroupService() *tradeUserGroupService {
	if TradeUserGroupService == nil {
		TradeUserGroupService = &tradeUserGroupService{
			redisRepo:          redis.InitRedisAndRedisClient(),
			eventBus:           getEventBusService(),
			firebaseDb:         getFirebaseRealTimeDatabase(),
			bullionService:     getBullionService(),
			tradeUserGroupRepo: repos.TradeUserRepo,
		}
		println("Trade User Group Service Initialized")
	}
	return TradeUserGroupService
}

// Create New Trade User Group And Create Mapping
func (t *tradeUserGroupService) CreateNewTradeUserGroup(bullionId string) {

}
