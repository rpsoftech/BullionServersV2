package services

import (
	"crypto/rand"
	"fmt"
	"io"

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

func (s *sendMsgService) SendOtp(otpReq *interfaces.OTPReqBase) error {
	data := s.redisRepo.GetStringData("otp/" + otpReq.BullionId + "/" + otpReq.Number)
	if len(data) > 0 {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_OTP_ALREADY_SENT,
			Message:    "Otp Already Sent",
			Name:       "ERROR_OTP_ALREADY_SENT",
		}
	}
	msgTemplate := new(interfaces.MsgTemplateBase)
	s.firebaseDb.GetData("msgTemplates", []string{otpReq.BullionId, "otp"}, msgTemplate)

	fmt.Printf("%#v", msgTemplate)
	println("here")
	// if err != nil {
	// 	fmt.Printf("%#v", err)
	// } else {
	// }
	return nil
}

// func (s *sendMsgService)
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateOTP(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
