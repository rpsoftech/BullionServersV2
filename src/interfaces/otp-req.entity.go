package interfaces

import "time"

type (
	OTPReqBaseEntity struct {
		BullionId string    `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Number    string    `bson:"number" json:"number" validate:"required,min=10,max=12"`
		Attempt   int16     `bson:"attempt" json:"attempt" validate:"required"`
		ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt" validate:"required"`
	}
)
