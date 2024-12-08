package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type bankRateService struct {
	bankRateRepo            *repos.BankRateCalcRepoStruct
	eventBus                *eventBusService
	firebaseDatabaseService *firebaseDatabaseService
}

var BankRateCalcService *bankRateService

func init() {
	getBankRateService()
}

func getBankRateService() *bankRateService {
	if BankRateCalcService == nil {
		BankRateCalcService = &bankRateService{
			eventBus:                getEventBusService(),
			bankRateRepo:            repos.BankRateCalcRepo,
			firebaseDatabaseService: getFirebaseRealTimeDatabase(),
		}
		println("Bank Rate Service Initialized")
	}
	return BankRateCalcService
}

func (service *bankRateService) GetBankRateCalcByBullionId(bullionId string) (*interfaces.BankRateCalcEntity, error) {
	return service.bankRateRepo.FindOneByBullionId(bullionId)
}

func (service *bankRateService) SaveBankRateCalc(gold *interfaces.BankRateCalcBase, silver *interfaces.BankRateCalcBase, bullionId string, adminId string) (*interfaces.BankRateCalcEntity, error) {
	entity, err := service.GetBankRateCalcByBullionId(bullionId)
	if err != nil {
		entity = &interfaces.BankRateCalcEntity{
			BullionId:   bullionId,
			GOLD_SPOT:   gold,
			SILVER_SPOT: silver,
		}
		entity.CreateNewBankRateCalc()
	} else {
		entity.GOLD_SPOT = gold
		entity.SILVER_SPOT = silver
	}
	_, err = service.bankRateRepo.Save(entity)
	if err != nil {
		return nil, err
	}
	service.eventBus.Publish(events.BankRateCalcUpdatedEvent(entity, adminId))
	service.firebaseDatabaseService.SetPublicData(bullionId, []string{"bankRateCalc"}, entity)
	return entity, nil
}
