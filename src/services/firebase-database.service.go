package services

import (
	"firebase.google.com/go/v4/db"
	"github.com/rpsoftech/bullion-server/src/utility/firebase"
)

type firebaseDatabaseService struct {
	db *db.Client
}

var FirebaseDatabaseService *firebaseDatabaseService

func init() {
	FirebaseDatabaseService = &firebaseDatabaseService{
		db: firebase.FirebaseDb,
	}
}

func (s *firebaseDatabaseService) SetData(bullionId string, path []string, data interface{}) error {
	ref := s.db.NewRef("bullions/" + bullionId)

	for _, child := range path {
		ref = ref.Child(child)
	}
	return ref.Set(firebase.FirebaseCtx, data)
}
