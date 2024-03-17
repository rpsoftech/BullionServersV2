package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type orderGeneralService struct {
	eventBus            *eventBusService
	firebaseDb          *firebaseDatabaseService
	bullionService      *bullionDetailsService
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
			orderRepo:           repos.OrderRepo,
			productGroupMapRepo: repos.ProductGroupMapRepo,
		}
		println("Order General Service Initialized")
	}
	return OrderGeneralService
}

func (service *orderGeneralService) ValidateVolumeForGroupMapId(groupMapId string, weight int) (bool, error) {
	groupMap, err := service.productGroupMapRepo.FindOne(groupMapId)
	if err != nil {
		return false, err
	}
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

func (service *orderGeneralService) ValidationOfGroupMapIdAndOrder(groupMapId string, userId string, weight int) (bool, error) {
	groupMap, err := service.productGroupMapRepo.FindOne(groupMapId)

	if err != nil {
		return false, err
	}

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
