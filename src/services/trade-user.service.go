package services

import "github.com/rpsoftech/bullion-server/src/mongodb/repos"

type tradeUserServiceStruct struct {
	tradeUserRepo    *repos.TradeUserRepoStruct
	eventBus         *eventBusService
	realtimeDatabase *firebaseDatabaseService
}

var TradeUserService *tradeUserServiceStruct

func init() {
	TradeUserService = &tradeUserServiceStruct{
		tradeUserRepo:    repos.TradeUserRepo,
		eventBus:         getEventBusService(),
		realtimeDatabase: getFirebaseRealTimeDatabase(),
	}
	println("Trade Service Initialized")
}
