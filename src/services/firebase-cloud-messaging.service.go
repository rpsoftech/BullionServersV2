package services

import (
	"firebase.google.com/go/v4/messaging"
	"github.com/rpsoftech/bullion-server/src/utility/firebase"
)

type firebaseCloudMessagingService struct {
	fcm *messaging.Client
}

var FirebaseCloudMessagingService *firebaseCloudMessagingService

func init() {
	FirebaseCloudMessagingService = &firebaseCloudMessagingService{
		fcm: firebase.FirebaseFCM,
	}
}

// func (s *firebaseCloudMessagingService) GenerateCustomToken(uid string, claims map[string]interface{}) (string, error) {
// 	token, err := s.auth.CustomTokenWithClaims(firebase.FirebaseCtx, uid, claims)
// 	if err != nil {
// 		err = &interfaces.RequestError{
// 			StatusCode: 500, Code: interfaces.ERROR_INTERNAL_SERVER, Message: "Issue In Generating Firebase Token", Name: "INTERNAL_SERVER_ERROR", Extra: err,
// 		}
// 		return token, err
// 	}
// 	return token, err
// }
