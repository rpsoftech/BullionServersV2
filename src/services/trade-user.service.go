package services

import (
	"time"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type tradeUserServiceStruct struct {
	tradeUserRepo    *repos.TradeUserRepoStruct
	eventBus         *eventBusService
	bullionService   *bullionDetailsService
	sendMsgService   *sendMsgService
	realtimeDatabase *firebaseDatabaseService
}

var TradeUserService *tradeUserServiceStruct

func init() {
	TradeUserService = &tradeUserServiceStruct{
		tradeUserRepo:    repos.TradeUserRepo,
		eventBus:         getEventBusService(),
		sendMsgService:   getSendMsgService(),
		bullionService:   getBullionService(),
		realtimeDatabase: getFirebaseRealTimeDatabase(),
	}
}

func (service *tradeUserServiceStruct) VerifyAndSendOtpForNewUser(tradeUser *interfaces.TradeUserBase, bullionId string) (string, error) {
	users, err := service.tradeUserRepo.FindDuplicateUser(tradeUser.Email, tradeUser.Number, tradeUser.Email, bullionId)
	if err != nil {
		return "", err
	}
	if len(*users) > 0 {
		return "", &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_DUPLICATE_USER,
			Message:    "User Exists With Matching With Wither Email,Number Or Username",
			Name:       "ERROR_DUPLICATE_USER",
		}
	}
	_, err = service.SendOtp(tradeUser.Name, tradeUser.Number, tradeUser.BullionId)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (service *tradeUserServiceStruct) SendOtp(name string, number string, bullionId string) (*interfaces.OTPReqEntity, error) {
	bullionDetails, err := service.bullionService.GetBullionDetailsByBullionId(bullionId)
	if err != nil {
		return nil, err
	}
	entity, err := service.sendMsgService.SendOtp(&interfaces.OTPReqBase{
		BullionId: bullionId,
		Number:    number,
		Attempt:   0,
		ExpiresAt: time.Now(),
	}, &interfaces.OTPReqVariablesStruct{
		BullionName: bullionDetails.Name,
		Name:        name,
		Number:      number,
	}, bullionDetails.BullionConfigs.OTPLength)

	if err != nil {
		return entity, err
	}
	return entity, nil
}