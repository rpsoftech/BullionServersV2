package events

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type bankRateCalcEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *bankRateCalcEvent) Add() *bankRateCalcEvent {
	base.ParentNames = []string{base.EventName, "BankRateCalcEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func BankRateCalcUpdatedEvent(entity *interfaces.BankRateCalcEntity, adminId string) *BaseEvent {
	event := &bankRateCalcEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "BankRateCalcUpdatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}
