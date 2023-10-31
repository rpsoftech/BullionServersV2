package events

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type otpReqEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *otpReqEvent) Add() *otpReqEvent {
	base.ParentNames = []string{base.EventName, "OtpReqEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateOtpSentEvent(entity *interfaces.OTPReqEntity, adminId string) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "OTPSent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateOtpResendEvent(entity *interfaces.OTPReqEntity, adminId string) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "OTPResent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateOtpVerifiedEvent(entity *interfaces.OTPReqEntity, adminId string) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			AdminId:   adminId,
			Payload:   entity,
			EventName: "OTPVerified",
		},
	}
	event.Add()
	return event.BaseEvent
}
