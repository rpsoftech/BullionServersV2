package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/redis"
)

type tradeUserGroupService struct {
	redisRepo           *redis.RedisClientStruct
	eventBus            *eventBusService
	firebaseDb          *firebaseDatabaseService
	bullionService      *bullionDetailsService
	tradeUserGroupRepo  *repos.TradeUserGroupRepoStruct
	productService      *productService
	productGroupMapRepo *repos.ProductGroupMapRepoStruct
}

var TradeUserGroupService *tradeUserGroupService

func init() {
	getTradeUserGroupService()
}

func getTradeUserGroupService() *tradeUserGroupService {
	if TradeUserGroupService == nil {
		TradeUserGroupService = &tradeUserGroupService{
			redisRepo:           redis.InitRedisAndRedisClient(),
			eventBus:            getEventBusService(),
			firebaseDb:          getFirebaseRealTimeDatabase(),
			bullionService:      getBullionService(),
			productService:      getProductService(),
			tradeUserGroupRepo:  repos.TradeUserGroupRepo,
			productGroupMapRepo: repos.ProductGroupMapRepo,
		}
		println("Trade User Group Service Initialized")
	}
	return TradeUserGroupService
}

// Create New Trade User Group And Create Mapping
func (t *tradeUserGroupService) CreateNewTradeUserGroup(bullionId string, name string) (*interfaces.TradeUserGroupEntity, error) {
	entity := &interfaces.TradeUserGroupEntity{
		BaseEntity: &interfaces.BaseEntity{},
		TradeUserGroupBase: &interfaces.TradeUserGroupBase{
			BullionId: bullionId,
			Name:      name,
			IsActive:  false,
			CanTrade:  false,
			CanLogin:  false,
		},
	}
	entity.CreateNew()
	if _, err := t.tradeUserGroupRepo.Save(entity); err != nil {
		return nil, err
	}
	err := t.createGroupMapFromNewGroup(entity.ID, bullionId)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (t *tradeUserGroupService) createGroupMapFromNewGroup(groupId string, bullionId string) error {
	entities, err := t.productService.GetProductsByBullionId(bullionId)
	if err != nil {
		return err
	}
	groupMapEntity := make([]interfaces.TradeUserGroupMapEntity, len(*entities))

	for i, entity := range *entities {
		groupMapEntity[i] = interfaces.TradeUserGroupMapEntity{
			BaseEntity: &interfaces.BaseEntity{},
			TradeUserGroupMapBase: &interfaces.TradeUserGroupMapBase{
				BullionId: bullionId,
				GroupId:   groupId,
				ProductId: entity.ID,
				IsActive:  false,
				CanTrade:  false,
				GroupPremiumBase: &interfaces.GroupPremiumBase{
					Buy:  0,
					Sell: 0,
				},
				GroupVolumeBase: &interfaces.GroupVolumeBase{
					OneClick: 0,
					Step:     0,
					Total:    0,
				},
			},
		}
		groupMapEntity[i].CreateNew()
	}
	t.productGroupMapRepo.BulkUpdate(&groupMapEntity)
	return nil
}
