package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
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

func (service *tradeUserServiceStruct) VerifyAndSendOtpForNewUser(tradeUser *interfaces.TradeUserBase, bullionId string) (*string, error) {
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
	tokenString, err := service.accessTokenService.GenerateToken(&localJwt.GeneralPurposeTokenGeneration{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Minute * 2)},
		},
		BullionId: bullionId,
		ExtraClaim: map[string]interface{}{
			"otpReqEntityId": otpReqEntity.ID,
			"tradeUser":      tradeUser,
		},
	})
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (service *tradeUserServiceStruct) verifyRegistrationToken(token string, returnTradeUser bool) (*localJwt.GeneralPurposeTokenGeneration, string, *interfaces.TradeUserBase, error) {
	claims, err := service.accessTokenService.VerifyTokenGeneralPurpose(token)
	if err != nil {
		return nil, "", nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "OTP Req Token Expired",
			Name:       "ERROR_INVALID_INPUT",
			Extra:      err,
		}
	}
	otpReqId, ok := claims.ExtraClaim["otpReqEntityId"].(string)
	if !ok {
		return nil, "", nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "OTP Req Id Not Found",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	tradeUserMap, ok := claims.ExtraClaim["tradeUser"]

	if !ok {
		return nil, "", nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "TradeUser Details Not Found",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	tradeUser := new(interfaces.TradeUserBase)
	err = mapstructure.Decode(tradeUserMap, &tradeUser)
	if err != nil {
		return nil, "", nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "TradeUser Details Not Found",
			Name:       "ERROR_INVALID_INPUT",
		}
	}
	return claims, otpReqId, tradeUser, nil
}

func (service *tradeUserServiceStruct) VerifyTokenAndResendOTP(token string) (*string, error) {
	claim, otpReqId, tradeUser, err := service.verifyRegistrationToken(token, true)
	fmt.Printf("Trade User %#v \n", tradeUser)
	if err != nil {
		return nil, err
	}
	err = service.sendMsgService.ResendOtp(otpReqId)
	if err != nil {
		return nil, err
	}
	claim.ExpiresAt = &jwt.NumericDate{Time: time.Now().Add(time.Minute * 2)}
	tokenString, err := service.accessTokenService.GenerateToken(claim)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (service *tradeUserServiceStruct) SendOtp(name string, number string, bullionId string) (*interfaces.OTPReqEntity, error) {
	bullionDetails, err := service.bullionService.GetBullionDetailsByBullionId(bullionId)
	if err != nil {
		return nil, err
	}
	entity, err := service.sendMsgService.SendOtp(&interfaces.OTPReqBase{
		BullionId: bullionId,
		Number:    number,
		Name:      name,
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

func (service *tradeUserServiceStruct) VerifyTokenAndVerifyOTP() {}
