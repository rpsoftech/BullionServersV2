package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/redis"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type sendMsgService struct {
	redisRepo                *redis.RedisClientStruct
	eventBus                 *eventBusService
	firebaseDb               *firebaseDatabaseService
	bullionService           *bullionDetailsService
	regExpForMessageVariable *regexp.Regexp
}

var SendMsgService *sendMsgService

func getSendMsgService() *sendMsgService {
	if SendMsgService == nil {
		SendMsgService = &sendMsgService{
			redisRepo:                redis.InitRedisAndRedisClient(),
			eventBus:                 getEventBusService(),
			firebaseDb:               getFirebaseRealTimeDatabase(),
			bullionService:           getBullionService(),
			regExpForMessageVariable: regexp.MustCompile(`##\S*##`),
		}
		println("Send Msg Service Initialized")
	}
	return SendMsgService
}

func (s *sendMsgService) SendOtp(otpReq *interfaces.OTPReqBase, variable *interfaces.OTPReqVariablesStruct, otpLength int) (*interfaces.OTPReqEntity, error) {
	data := s.redisRepo.GetStringData("otp/" + otpReq.BullionId + "/" + otpReq.Number)
	if len(data) > 0 {
		return nil, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_OTP_ALREADY_SENT,
			Message:    "Otp Already Sent",
			Name:       "ERROR_OTP_ALREADY_SENT",
		}
	}
	variable.OTP = GenerateOTP(otpLength)
	entity := interfaces.CreateOTPEntity(otpReq, variable.OTP)
	err := s.prepareAndSendOTP(entity, variable)
	if err != nil {
		return nil, err
	}
	s.eventBus.Publish(events.CreateOtpSentEvent(entity))
	return entity, nil
}

func (s *sendMsgService) ResendOtp(otpReqId string) error {
	data := s.redisRepo.GetStringData("otp/" + otpReqId)
	if data == "" {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_OTP_EXPIRED,
			Message:    "OTP Req Expired",
			Name:       "ERROR_OTP_EXPIRED",
		}
	}
	otpReqEntity := new(interfaces.OTPReqEntity)
	err := json.Unmarshal([]byte(data), otpReqEntity)
	if err != nil {
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_INTERNAL_SERVER,
			Message:    "Unable to parse OTP REQ JSON",
		}
	}
	otpReqEntity.RestoreTimeStamp()
	if time.Now().Before(otpReqEntity.ModifiedAt.Add(time.Second * 15)) {
		return &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_OTP_ALREADY_SENT,
			Message:    "Please Wait For 10 Seconds Before Requesting",
			Name:       "ERROR_OTP_ALREADY_SENT",
		}
	}
	otpReqEntity.NewAttempt()
	bullionDetails, err := s.bullionService.GetBullionDetailsByBullionId(otpReqEntity.BullionId)
	if err != nil {
		return err
	}
	err = s.prepareAndSendOTP(otpReqEntity, &interfaces.OTPReqVariablesStruct{
		OTP:         otpReqEntity.OTP,
		Name:        otpReqEntity.Name,
		Number:      otpReqEntity.Number,
		BullionName: bullionDetails.Name,
	})
	if err != nil {
		return err
	}
	s.eventBus.Publish(events.CreateOtpResendEvent(otpReqEntity))
	return err
}

func (s *sendMsgService) prepareAndSendOTP(otpReq *interfaces.OTPReqEntity, variable *interfaces.OTPReqVariablesStruct) error {
	msgTemplate := new(interfaces.MsgTemplateBase)
	err := s.firebaseDb.GetData("msgTemplates", []string{otpReq.BullionId, "otp"}, msgTemplate)
	if msgTemplate.WhatsappTemplate == "" && msgTemplate.MSG91Id == "" {
		// TODO Throw Critical Error Which needs to be reported
		println("Something Went Wrong While Fetching OTP Templates")
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_WHILE_FETCHING_MESSAGE_TEMPLATE,
			Message:    "OTP Template Error",
			Name:       "ERROR_WHILE_FETCHING_MESSAGE_TEMPLATE",
			Extra:      msgTemplate,
		}
	}
	if err != nil {
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_WHILE_FETCHING_MESSAGE_TEMPLATE,
			Message:    "OTP Template NOT Found",
			Name:       "ERROR_WHILE_FETCHING_MESSAGE_TEMPLATE",
			Extra:      err,
		}
	}
	err = s.saveAndUpdateOTPService(otpReq)
	if err != nil {
		return err
	}
	err = s.sendWhatsappMessage(msgTemplate.WhatsappTemplate, "OTP", variable, otpReq)
	if err != nil {
		return err
	}
	return nil
}

func (s *sendMsgService) saveAndUpdateOTPService(otpEntity *interfaces.OTPReqEntity) error {
	otpEntity.ExpiresAt = otpEntity.ExpiresAt.Add(120 * time.Second)
	otpEntity.AddTimeStamps()
	fmt.Printf("%#v", otpEntity)
	otpEntityStringBytes, err := json.Marshal(otpEntity)
	if err != nil {
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_INTERNAL_SERVER,
			Message:    "Unable convert OTP REQ to string",
			Name:       "OTPReq Entity Marshal Error",
			Extra:      err,
		}
	}
	otpEntityString := string(otpEntityStringBytes)
	s.redisRepo.SetStringData(fmt.Sprintf("otp/%s/%s", otpEntity.BullionId, otpEntity.Number), otpEntityString, 120)
	s.redisRepo.SetStringData(fmt.Sprintf("otp/%s", otpEntity.ID), otpEntityString, 120)
	return nil
}

func (s *sendMsgService) sendWhatsappMessage(template string, templateName string, variables interface{}, otpReqEntity *interfaces.OTPReqEntity) error {
	jsonMap, err := utility.StructToStringMap(variables)
	if err != nil {
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_INTERNAL_SERVER,
			Message:    "Error While converting Struct To JSON",
			Name:       "CONVERSATION_ERROR",
			Extra:      err,
		}
	}
	routeToPost := otpReqEntity.BullionId
	bullionDetails, err := s.bullionService.GetBullionDetailsByBullionId(otpReqEntity.BullionId)
	if err != nil {
		return err
	}
	if !bullionDetails.BullionConfigs.HaveCustomWhatsappAgent {
		routeToPost = "common"
	}
	message := s.processMessage(template, &jsonMap)
	err = s.firebaseDb.setPrivateData("whatsappMessage", []string{routeToPost, otpReqEntity.ID}, map[string]string{
		"message": message,
		"number":  otpReqEntity.Number,
	})
	if err != nil {
		return &interfaces.RequestError{
			StatusCode: http.StatusInternalServerError,
			Code:       interfaces.ERROR_INTERNAL_SERVER,
			Message:    "Error Posting Whatsapp Message to Firebase",
			Name:       "CONVERSATION_ERROR",
			Extra:      err,
		}
	}
	s.eventBus.Publish(events.CreateWhatsappMessageSendEvent(otpReqEntity.BullionId, templateName, otpReqEntity.Number, message))
	return nil
}
func (s *sendMsgService) processMessage(template string, variables *map[string]string) string {
	for key, value := range *variables {
		template = strings.ReplaceAll(template, fmt.Sprintf("##%s##", key), value)
	}
	template = s.regExpForMessageVariable.ReplaceAllString(template, "")
	return template
}

// func (s *sendMsgService) sendMsg91(){}

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
