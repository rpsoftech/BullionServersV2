package services

import "github.com/rpsoftech/bullion-server/src/mongodb/repos"

type orderGeneralService struct {
	eventBus       *eventBusService
	firebaseDb     *firebaseDatabaseService
	bullionService *bullionDetailsService
	orderRepo      *repos.OrderRepoStruct
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
			orderRepo:      repos.OrderRepo,
		}
		println("Order General Service Initialized")
	}
	return OrderGeneralService
}
