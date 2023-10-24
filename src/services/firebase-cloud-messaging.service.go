package services

import (
	"strconv"

	"firebase.google.com/go/v4/messaging"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/utility/firebase"
)

type firebaseCloudMessagingService struct {
	fcm *messaging.Client
}

var fcmService *firebaseCloudMessagingService

func getFirebaseCloudMessagingService() *firebaseCloudMessagingService {
	if fcmService == nil {
		fcmService = &firebaseCloudMessagingService{
			fcm: firebase.FirebaseFCM,
		}
		println("FCM Service Initialized")
	}
	return fcmService
}

func (s *firebaseCloudMessagingService) SendTextNotificationToAll(bullionId string, title string, body string, isHtml bool) {
	s.SendToTopic(bullionId, messaging.Notification{
		Title: title,
		Body:  body,
	}, map[string]string{
		"title":  title,
		"body":   body,
		"isHtml": strconv.FormatBool(isHtml),
	})
}

func (s *firebaseCloudMessagingService) SendToTopic(bullionId string, notification messaging.Notification, extra map[string]string) {
	// TODO: NEED To Add TTL
	s.fcm.Send(firebase.FirebaseCtx, &messaging.Message{
		Data:         extra,
		Notification: &notification,
		Topic:        bullionId + "/main",
	})
}

func (s *firebaseCloudMessagingService) SubscribeToChanel(bullionId string, token string, deviceType interfaces.DeviceType) {
	s.fcm.SubscribeToTopic(firebase.FirebaseCtx, []string{token}, bullionId+"/main")
	s.fcm.SubscribeToTopic(firebase.FirebaseCtx, []string{token}, bullionId+"/"+deviceType.String())
}
