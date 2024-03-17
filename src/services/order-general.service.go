package services

import (
	"net/http"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type orderGeneralService struct {
	eventBus            *eventBusService
	firebaseDb          *firebaseDatabaseService
	bullionService      *bullionDetailsService
	flagService         *FlagServiceStruct
	groupMapService     *tradeUserGroupServiceStruct
	orderRepo           *repos.OrderRepoStruct
	productGroupMapRepo *repos.ProductGroupMapRepoStruct
}

var OrderGeneralService *orderGeneralService

func init() {
	getOrderGeneralService()
}
func getOrderGeneralService() *orderGeneralService {
	if OrderGeneralService == nil {
		OrderGeneralService = &orderGeneralService{
			eventBus:            getEventBusService(),
			firebaseDb:          getFirebaseRealTimeDatabase(),
			bullionService:      getBullionService(),
			groupMapService:     getTradeUserGroupService(),
			flagService:         getFlagService(),
			orderRepo:           repos.OrderRepo,
			productGroupMapRepo: repos.ProductGroupMapRepo,
		}
		println("Order General Service Initialized")
	}
	return OrderGeneralService
}

func (service *orderGeneralService) ValidateUserAndGroupMapWithWeight(user *interfaces.TradeUserEntity, groupMap *interfaces.TradeUserGroupMapEntity, group *interfaces.TradeUserGroupEntity, weight int) (bool, error) {
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
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED,
			Message:    "Trading is disabled. Contact User",
			Name:       "BULLION_NOT_ACTIVE",
		}
	}
	if !group.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

	if !group.CanTrade {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED_FOR_GROUP,
			Message:    "Trading is disabled for your group. Contact User",
			Name:       "GROUP_NOT_ACTIVE",
		}
	}

	if !groupMap.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Map Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

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
			Name:       "INVALID_REQUEST",
		}
	}
	return true, nil
}
