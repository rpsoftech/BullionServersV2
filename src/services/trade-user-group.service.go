package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
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
func (t *tradeUserGroupService) CreateNewTradeUserGroup(bullionId string, name string, adminId string) (*interfaces.TradeUserGroupEntity, error) {
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
	err := t.createGroupMapFromNewGroup(entity.ID, bullionId, adminId)
	if err != nil {
		return nil, err
	}
	if siteDetails, _ := t.bullionService.GetBullionDetailsByBullionId(bullionId); siteDetails != nil {
		if siteDetails.BullionConfigs.DefaultGroupIdForTradeUser == "" {
			siteDetails.BullionConfigs.DefaultGroupIdForTradeUser = entity.ID
			t.bullionService.UpdateBullionSiteDetails(siteDetails)
		}
	}
	t.eventBus.Publish(events.CreateTradeUserGroupCreated(bullionId, entity, adminId))
	return entity, nil
}

func (t *tradeUserGroupService) createGroupMapFromNewGroup(groupId string, bullionId string, adminId string) error {
	entities, err := t.productService.GetProductsByBullionId(bullionId)
	if err != nil {
		return err
	}
	groupMapEntities := make([]interfaces.TradeUserGroupMapEntity, len(*entities))

	for i, entity := range *entities {
		groupMapEntities[i] = interfaces.TradeUserGroupMapEntity{
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
		groupMapEntities[i].CreateNew()
	}
	t.productGroupMapRepo.BulkUpdate(&groupMapEntities)
	t.eventBus.Publish(events.CreateTradeUserGroupMapUpdated(bullionId, &groupMapEntities, groupId, adminId))
	return nil
}

func (t *tradeUserGroupService) CreateGroupMapFromNewProduct(productId string, bullionId string, adminId string) error {
	entities, err := t.tradeUserGroupRepo.GetAllByBullionId(bullionId)
	if err != nil {
		return err
	}
	groupMapEntities := make([]interfaces.TradeUserGroupMapEntity, len(*entities))
	for i, entity := range *entities {
		groupMapEntities[i] = interfaces.TradeUserGroupMapEntity{
			BaseEntity: &interfaces.BaseEntity{},
			TradeUserGroupMapBase: &interfaces.TradeUserGroupMapBase{
				BullionId: bullionId,
				GroupId:   entity.ID,
				ProductId: productId,
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
		groupMapEntities[i].CreateNew()
	}
	t.productGroupMapRepo.BulkUpdate(&groupMapEntities)
	go func() {
		for _, entity := range groupMapEntities {
			t.eventBus.Publish(events.CreateTradeUserGroupMapUpdated(bullionId, &[]interfaces.TradeUserGroupMapEntity{entity}, entity.GroupId, adminId))
		}
	}()
	return nil
}
