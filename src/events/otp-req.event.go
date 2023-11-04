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

func CreateOtpSentEvent(entity *interfaces.OTPReqEntity) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			Payload:   entity,
			EventName: "OTPSent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateOtpResendEvent(entity *interfaces.OTPReqEntity) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			Payload:   entity,
			EventName: "OTPResent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateOtpVerifiedEvent(entity *interfaces.OTPReqEntity) *BaseEvent {
	event := &otpReqEvent{
		BaseEvent: &BaseEvent{
			BullionId: entity.BullionId,
			KeyId:     entity.ID,
			Payload:   entity,
			EventName: "OTPVerified",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateWhatsappMessageSendEvent(bullionId string, templateId string, number string, message string) *BaseEvent {
	event := &BaseEvent{
		BullionId: bullionId,
		KeyId:     templateId,
		Payload: map[string]string{
			"number":  number,
			"message": message,
		},
		ParentNames: []string{"WhatsappMessageSend", "MessageEvent"},
		EventName:   "WhatsappMessageSend",
	}
	event.CreateBaseEvent()
	return event
}
