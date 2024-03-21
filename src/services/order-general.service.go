package services

import (
	"net/http"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type orderGeneralService struct {
	eventBus       *eventBusService
	firebaseDb     *firebaseDatabaseService
	bullionService *bullionDetailsService
	flagService    *FlagServiceStruct
	groupService   *tradeUserGroupServiceStruct
	orderRepo      *repos.OrderRepoStruct
	userService    *tradeUserServiceStruct
	productService *productService
}

var OrderGeneralService *orderGeneralService

func init() {
	getOrderGeneralService()
}
func getOrderGeneralService() *orderGeneralService {
	if OrderGeneralService == nil {
		OrderGeneralService = &orderGeneralService{
			eventBus:       getEventBusService(),
			firebaseDb:     getFirebaseRealTimeDatabase(),
			bullionService: getBullionService(),
			groupService:   getTradeUserGroupService(),
			flagService:    getFlagService(),
			userService:    getTradeUserService(),
			productService: getProductService(),
			orderRepo:      repos.OrderRepo,
		}
		println("Order General Service Initialized")
	}
	return OrderGeneralService
}

func (service *orderGeneralService) ValidateUserAndGroupMapWithWeight(user *interfaces.TradeUserEntity, group *interfaces.TradeUserGroupEntity, groupMap *interfaces.TradeUserGroupMapEntity, weight int) (bool, error) {

	// Check for User Activation
	if !user.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Account Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

	if flags, err := service.flagService.GetFlags(user.BullionId); err != nil {
		return false, err
	} else if !flags.CanTrade {
		// Check If Trading Is Disabled
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED,
			Message:    "Trading is disabled. Contact User",
			Name:       "BULLION_NOT_ACTIVE",
		}
	}
	// Check for Group Activation
	if !group.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}
	// Check for Group User Can Trade
	if !group.CanTrade {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED_FOR_GROUP,
			Message:    "Trading is disabled for your group. Contact User",
			Name:       "GROUP_NOT_ACTIVE",
		}
	}

	// Check for Group Map Activation
	if !groupMap.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Map Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

	// Check for Group Map Can Trade
	if !groupMap.CanTrade {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED_FOR_PRODUCT,
			Message:    "Trading is disabled for your group map. Contact Admin",
			Name:       "GROUP_NOT_ACTIVE",
		}
	}

	return service.validateVolumeForGroupMap(groupMap, weight)
}

func (service *orderGeneralService) validateVolumeForGroupMap(groupMap *interfaces.TradeUserGroupMapEntity, weight int) (bool, error) {
	if !groupMap.ValidateVolume(weight) {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_VOLUME,
			Message:    "Invalid Volume",
			Name:       "INVALID_VOLUME",
		}
	}
	return true, nil
}

func (service *orderGeneralService) findOrderDetailsAndValidate(userId string, groupId string, groupMapId string, weight int) (*interfaces.TradeUserEntity, *interfaces.TradeUserGroupEntity, *interfaces.TradeUserGroupMapEntity, error) {
	// Get User
	user, err := service.userService.GetTradeUserById(userId)
	if err != nil {
		return nil, nil, nil, err
	}

	if user.GroupId != groupId {
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "MissMatch Group Id",
			Name:       "MISS_MATCH_GROUP_ID",
		}
	}
	// Get Group
	group, err := service.groupService.GetGroupByGroupId(groupId, user.BullionId)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get Group Map
	groupMaps, err := service.groupService.GetGroupMapByGroupId(groupId, user.BullionId)
	if err != nil {
		return nil, nil, nil, err
	}
	var groupMap *interfaces.TradeUserGroupMapEntity

	for _, v := range *groupMaps {
		if v.ID == groupMapId {
			groupMap = &v
			break
		}
	}
	if groupMap == nil {
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GROUP_MAP_NOT_FOUND,
			Message:    "Group Map Not Found",
			Name:       "GROUP_MAP_NOT_FOUND",
		}

	}
	service.ValidateUserAndGroupMapWithWeight(user, group, groupMap, weight)
	return user, group, groupMap, nil
}
func (service *orderGeneralService) PlaceOrder(orderType interfaces.OrderStatus, userId string, groupId string, groupMapId string, weight int, price float64, placedBy string) (*interfaces.OrderEntity, error) {

	user, group, groupMap, err := service.findOrderDetailsAndValidate(userId, groupId, groupMapId, weight)
	if err != nil {
		return nil, err
	}
	// product
	product, err := service.productService.GetProductsById(group.BullionId, groupMap.ProductId)
	if err != nil {
		return nil, err
	}

	// TODO Validate Pricing

	_, err = user.UpdateMarginAfterOrder(weight, product.SourceSymbol)
	if err != nil {
		return nil, err
	}

	// TODO Check Hedging And Place Order
	return nil, nil
	// return service.orderRepo.PlaceOrder(orderType, user, group, groupMap, price, placedBy)
	// return service.orderRepo.PlaceOrder(orderType, user, group, groupMap, price, placedBy)
}
