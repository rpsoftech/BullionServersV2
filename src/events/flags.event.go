package events

import "github.com/rpsoftech/bullion-server/src/interfaces"

type flagsEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *flagsEvent) Add() *flagsEvent {
	base.ParentNames = []string{base.EventName, "FlagsEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func FlagsUpdatedEvent(entity *interfaces.FlagsInterface, adminId string) *BaseEvent {
	event := &flagsEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.BullionId,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "FlagsUpdatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}
