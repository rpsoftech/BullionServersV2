package events

import (
	"github.com/google/uuid"
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type feedEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *feedEvent) Add() *feedEvent {
	base.ParentNames = []string{base.EventName, "FeedEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateUpdateFeedEvent(entity *interfaces.FeedsEntity, adminId string) *BaseEvent {
	event := &feedEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "FeedCreatedUpdatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateDeleteFeedEvent(entity *interfaces.FeedsBase, id string, adminId string) *BaseEvent {
	event := &feedEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     id,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "FeedDeletedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateNotificationSendEvent(entity *interfaces.FeedsBase, adminId string) *BaseEvent {
	event := &feedEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     uuid.New().String(),
			AdminId:   adminId,
			Payload:   entity,
			EventName: "NotificationSendEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}
