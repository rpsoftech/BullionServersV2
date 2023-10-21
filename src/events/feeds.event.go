package events

import (
	"github.com/google/uuid"
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type feedEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *feedEvent) Add() *feedEvent {
	base.ParentNames = append(base.ParentNames, "FeedEvent")
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateUpdateFeedEvent(entity *interfaces.FeedsEntity, adminId string) *BaseEvent {
	event := &feedEvent{
		BaseEvent: &BaseEvent{
			BullionId:   entity.BullionId,
			KeyId:       entity.ID,
			AdminId:     adminId,
			Payload:     entity,
			EventName:   "FeedCreatedUpdatedEvent",
			ParentNames: []string{"FeedCreatedUpdatedEvent"},
		},
	}
	event.Add()
	return event.BaseEvent
}
func CreateNotificationSendEvent(entity *interfaces.FeedsBase, adminId string) *BaseEvent {
	event := &feedEvent{
		BaseEvent: &BaseEvent{
			BullionId:   entity.BullionId,
			KeyId:       uuid.New().String(),
			AdminId:     adminId,
			Payload:     entity,
			EventName:   "NotificationSendEvent",
			ParentNames: []string{"NotificationSendEvent"},
		},
	}
	event.Add()
	return event.BaseEvent
}
