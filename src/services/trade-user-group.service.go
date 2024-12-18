package services

import (
	"net/http"

	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type tradeUserGroupServiceStruct struct {
	eventBus                      *eventBusService
	firebaseDb                    *firebaseDatabaseService
	bullionService                *bullionDetailsService
	tradeUserGroupRepo            *repos.TradeUserGroupRepoStruct
	productService                *productService
	productGroupMapRepo           *repos.ProductGroupMapRepoStruct
	groupMapByGroupIdMap          map[string]*[]interfaces.TradeUserGroupMapEntity
	groupsByBullionIdMapStructure map[string]*[]interfaces.TradeUserGroupEntity
	groupByGroupIdMapStructure    map[string]*interfaces.TradeUserGroupEntity
}

var TradeUserGroupService *tradeUserGroupServiceStruct

func init() {
	getTradeUserGroupService()
}

func getTradeUserGroupService() *tradeUserGroupServiceStruct {
	if TradeUserGroupService == nil {
		TradeUserGroupService = &tradeUserGroupServiceStruct{
			eventBus:                      getEventBusService(),
			firebaseDb:                    getFirebaseRealTimeDatabase(),
			bullionService:                getBullionService(),
			productService:                getProductService(),
			tradeUserGroupRepo:            repos.TradeUserGroupRepo,
			productGroupMapRepo:           repos.ProductGroupMapRepo,
			groupMapByGroupIdMap:          map[string]*[]interfaces.TradeUserGroupMapEntity{},
			groupsByBullionIdMapStructure: map[string]*[]interfaces.TradeUserGroupEntity{},
			groupByGroupIdMapStructure:    map[string]*interfaces.TradeUserGroupEntity{},
		}
		println("Trade User Group Service Initialized")
	}
	return TradeUserGroupService
}

// Create New Trade User Group And Create Mapping
func (t *tradeUserGroupServiceStruct) CreateNewTradeUserGroup(bullionId string, name string, adminId string) (*interfaces.TradeUserGroupEntity, error) {
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
	delete(t.groupsByBullionIdMapStructure, bullionId)
	t.eventBus.Publish(events.CreateTradeUserGroupCreated(bullionId, entity, adminId))
	t.updateGroupInFirebase(bullionId, entity)
	return entity, nil
}

func (t *tradeUserGroupServiceStruct) createGroupMapFromNewGroup(groupId string, bullionId string, adminId string) error {
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
	t.updateGroupMapInFirebase(bullionId, groupId, &groupMapEntities)
	t.eventBus.Publish(events.CreateTradeUserGroupMapUpdated(bullionId, &groupMapEntities, groupId, adminId))
	return nil
}

func (t *tradeUserGroupServiceStruct) UpdateTradeGroup(base *interfaces.TradeUserGroupBase, groupId string, adminId string) (*interfaces.TradeUserGroupEntity, error) {
	entity, err := t.GetGroupByGroupId(groupId, base.BullionId)
	if err != nil {
		return nil, err
	}
	entity.TradeUserGroupBase = base

	if _, err := t.tradeUserGroupRepo.Save(entity); err != nil {
		return nil, err
	}
	delete(t.groupMapByGroupIdMap, groupId)
	delete(t.groupsByBullionIdMapStructure, base.BullionId)
	t.eventBus.Publish(events.CreateTradeUserGroupUpdated(base.BullionId, entity, adminId))
	go t.updateGroupInFirebase(base.BullionId, entity)
	return entity, nil
}

func (t *tradeUserGroupServiceStruct) UpdateTradeGroupMap(base *[]interfaces.TradeUserGroupMapEntity, groupId string, bullionId string, adminId string) (*[]interfaces.TradeUserGroupMapEntity, error) {
	entity, err := t.GetGroupMapByGroupId(groupId, bullionId)
	if err != nil {
		return nil, err
	}

	if len(*entity) != len(*base) {
		return nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "Please Pass All Group Map Entities",
			Name:       "PLEASE_PASS_ALL_GROUP_MAP_ENTITIES",
		}
	}
	for i, record := range *entity {
		for _, baseRecord := range *base {
			if record.ID == baseRecord.ID {
				(*entity)[i].TradeUserGroupMapBase.UpdateDetails(baseRecord.TradeUserGroupMapBase)
			}
		}
	}

	if _, err := t.productGroupMapRepo.BulkUpdate(entity); err != nil {
		return nil, err
	}
	// Clear Cache Of Product Group Map
	delete(t.groupMapByGroupIdMap, groupId)
	t.eventBus.Publish(events.CreateTradeUserGroupMapUpdated(bullionId, entity, groupId, adminId))
	return entity, nil
}

// Creating New Mapping Of Group And Product After Creating New Product
func (t *tradeUserGroupServiceStruct) CreateGroupMapFromNewProduct(productId string, bullionId string, adminId string) error {
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
			// Clearing Cache
			delete(t.groupMapByGroupIdMap, entity.GroupId)
			t.eventBus.Publish(events.CreateTradeUserGroupMapUpdated(bullionId, &[]interfaces.TradeUserGroupMapEntity{entity}, entity.GroupId, adminId))
		}
	}()
	return nil
}

func (t *tradeUserGroupServiceStruct) GetAllGroupsByBullionId(bullionId string) (*[]interfaces.TradeUserGroupEntity, error) {
	if entity, ok := t.groupsByBullionIdMapStructure[bullionId]; ok {
		return entity, nil
	}
	if entity, err := t.tradeUserGroupRepo.GetAllByBullionId(bullionId); err == nil || len(*entity) == 0 {
		t.groupsByBullionIdMapStructure[bullionId] = entity
		for _, group := range *entity {
			go t.updateGroupInFirebase(bullionId, &group)
		}
		return entity, nil
	}
	return nil, &interfaces.RequestError{
		StatusCode: http.StatusBadRequest,
		Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
		Message:    "Groups Not Found For This Bullion",
		Name:       "GROUPS_NOT_FOUND_FOR_THIS_BULLION",
	}
}

func (t *tradeUserGroupServiceStruct) GetGroupMapByGroupId(groupId string, bullionId string) (*[]interfaces.TradeUserGroupMapEntity, error) {
	if entity, ok := t.groupMapByGroupIdMap[groupId]; ok {
		return entity, nil
	}

	if entity, err := t.productGroupMapRepo.GetAllByGroupId(groupId, bullionId); err == nil || len(*entity) == 0 {
		t.groupMapByGroupIdMap[groupId] = entity
		go t.updateGroupMapInFirebase(bullionId, groupId, entity)
		return entity, nil
	}

	return nil, &interfaces.RequestError{
		StatusCode: http.StatusBadRequest,
		Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
		Message:    "Groups Map Not Found For This Bullion And Group Id",
		Name:       "GROUPS_MAP_NOT_FOUND_FOR_THIS_BULLION_AND_GROUP_ID",
	}
}

func (t *tradeUserGroupServiceStruct) GetGroupByGroupId(groupId string, bullionId string) (*interfaces.TradeUserGroupEntity, error) {
	if entity, ok := t.groupByGroupIdMapStructure[groupId]; ok {
		return entity, nil
	}
	if entity, err := t.tradeUserGroupRepo.FindOne(groupId); err == nil && entity.BullionId == bullionId {
		t.groupByGroupIdMapStructure[groupId] = entity
		go t.updateGroupInFirebase(bullionId, entity)
		return entity, nil
	}
	return nil, &interfaces.RequestError{
		StatusCode: http.StatusBadRequest,
		Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
		Message:    "Groups Map Not Found For This Bullion And Group Id",
		Name:       "GROUPS_MAP_NOT_FOUND_FOR_THIS_BULLION_AND_GROUP_ID",
	}
}

func (t *tradeUserGroupServiceStruct) updateGroupInFirebase(bullionId string, group *interfaces.TradeUserGroupEntity) error {
	return t.firebaseDb.SetPublicData(bullionId, []string{"trade", "groups", group.ID}, group)

}
func (t *tradeUserGroupServiceStruct) updateGroupMapInFirebase(bullionId string, groupId string, maps *[]interfaces.TradeUserGroupMapEntity) error {
	// return t.firebaseDb.SetPublicData(bullionId, []string{"trade", "groups", group.ID}, group)
	return t.firebaseDb.SetPublicData(bullionId, []string{"trade", "groupMaps", groupId}, maps)

}
