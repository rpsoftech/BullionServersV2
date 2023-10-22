package events

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type bankDetailsEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *bankDetailsEvent) Add() *bankDetailsEvent {
	base.ParentNames = []string{base.EventName, "BankDetailsEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateNewBankDetailsCreated(entity *interfaces.BankDetailsEntity, adminId string) *BaseEvent {
	event := &bankDetailsEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "BankDetailsCreatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateBankDetailsDeletedEvent(entity *interfaces.BankDetailsBase, id string, adminId string) *BaseEvent {
	event := &bankDetailsEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     id,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "BankDetailsDeletedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateBankDetailsUpdatedEvent(entity *interfaces.BankDetailsEntity, adminId string) *BaseEvent {
	event := &bankDetailsEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "BankDetailsUpdatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}
