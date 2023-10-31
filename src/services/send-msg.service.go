package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/redis"
)

type sendMsgService struct {
	redisRepo  *redis.RedisClientStruct
	eventBus   *eventBusService
	firebaseDb *firebaseDatabaseService
}

var SendMsgService *sendMsgService

func getSendMsgService() *sendMsgService {
	if SendMsgService == nil {
		SendMsgService = &sendMsgService{
			redisRepo:  redis.InitRedisAndRedisClient(),
			eventBus:   getEventBusService(),
			firebaseDb: getFirebaseRealTimeDatabase(),
		}
		println("Send Msg Service Initialized")
	}
	return SendMsgService
}

func (s *sendMsgService) SendOtp(otpReq *interfaces.OTPReqBase) {
	msgTemplate := new(interfaces.MsgTemplateBase)
	s.firebaseDb.GetData("msgTemplates", []string{otpReq.BullionId, "otp"}, msgTemplate)

	println("here")
	// if err != nil {
	// 	fmt.Printf("%#v", err)
	// } else {
	// 	fmt.Printf("%#v", msgTemplate)
	// }
}

// func (s *sendMsgService)
