package events

import (
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	ObjId       string      `bson:"_id" json:"_id"`
	Id          string      `bson:"id" json:"id"`
	BullionId   string      `bson:"bullionId" json:"bullionId"`
	KeyId       string      `bson:"key" json:"key"`
	EventName   string      `bson:"eventName" json:"eventName"`
	ParentNames []string    `bson:"parentNames" json:"parentNames"`
	Payload     interface{} `bson:"payload" json:"payload"`
	AdminId     string      `bson:"adminId" json:"adminId"`
	OccurredAt  time.Time   `bson:"occurredAt" json:"occurredAt"`
}

func (base *BaseEvent) CreateBaseEvent() *BaseEvent {
	base.Id = uuid.New().String()
	base.ObjId = base.Id
	base.OccurredAt = time.Now()
	return base
}
