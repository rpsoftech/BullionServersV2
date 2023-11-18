package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	localJwt "github.com/rpsoftech/bullion-server/src/utility/jwt"
)

type tradeUserServiceStruct struct {
	accessTokenService *localJwt.TokenService
	tradeUserRepo      *repos.TradeUserRepoStruct
	eventBus           *eventBusService
	bullionService     *bullionDetailsService
	sendMsgService     *sendMsgService
	realtimeDatabase   *firebaseDatabaseService
}

var TradeUserService *tradeUserServiceStruct

func init() {
	TradeUserService = &tradeUserServiceStruct{
		tradeUserRepo:      repos.TradeUserRepo,
		accessTokenService: AccessTokenService,
		eventBus:           getEventBusService(),
		sendMsgService:     getSendMsgService(),
		bullionService:     getBullionService(),
		realtimeDatabase:   getFirebaseRealTimeDatabase(),
	}
}

func (service *tradeUserServiceStruct) VerifyAndSendOtpForNewUser(tradeUser *interfaces.TradeUserBase, bullionId string) (*interfaces.ApiTradeUserRegisterResponse, error) {
	users, err := service.tradeUserRepo.FindDuplicateUser(tradeUser.Email, tradeUser.Number, tradeUser.Email, bullionId)
	if err != nil {
		return nil, err
	}
	if len(*users) > 0 {
		return nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_DUPLICATE_USER,
			Message:    "User Exists With Matching With Wither Email,Number Or Username",
			Name:       "ERROR_DUPLICATE_USER",
		}
	}
	otpReqEntity, err := service.SendOtp(tradeUser.Name, tradeUser.Number, tradeUser.BullionId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	otpReqEntityString, err := json.Marshal(otpReqEntity.OTPReqBase)
	if err != nil {
		return nil, err
	}
	otpReqToken, err := service.accessTokenService.GenerateToken(&localJwt.GeneralPurposeTokenGeneration{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Minute * 10)},
		},
		BullionId:  bullionId,
		ExtraClaim: string(otpReqEntityString),
	})
	if err != nil {
		return nil, err
	}
	tradeUserString, err := json.Marshal(tradeUser)
	if err != nil {
		return nil, err
	}
	tradeUserToken, err := service.accessTokenService.GenerateToken(&localJwt.GeneralPurposeTokenGeneration{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Minute * 10)},
		},
		BullionId:  bullionId,
		ExtraClaim: string(tradeUserString),
	})
	if err != nil {
		return nil, err
	}
	return &interfaces.ApiTradeUserRegisterResponse{
		UserToken:   tradeUserToken,
		OtpReqToken: otpReqToken,
	}, nil
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
