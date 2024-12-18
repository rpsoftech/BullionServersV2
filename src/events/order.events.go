package events

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type OrderEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *OrderEvent) Add() *OrderEvent {
	base.ParentNames = []string{base.EventName, "OrderEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func orderStatusUpdatedEvent(entity *interfaces.OrderEntity, adminId string, eventName string) *BaseEvent {
	event := &OrderEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: eventName,
		},
	}
	event.Add()
	return event.BaseEvent
}

func OrderPlacedEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "OrderPlacedEvent")
}

func LimitPlacedEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "LimitPlacedEvent")
}

func LimitDeletedEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "LimitDeletedEvent")
}

func LimitPassedEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "LimitPassedEvent")
}

func LimitCanceledEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "LimitCanceledEvent")
}

func OrderDeliveredEvent(entity *interfaces.OrderEntity, adminId string) *BaseEvent {
	return orderStatusUpdatedEvent(entity, adminId, "OrderDeliveredEvent")
}
