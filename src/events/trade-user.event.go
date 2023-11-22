package events

import "github.com/rpsoftech/bullion-server/src/interfaces"

type tradeUserEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *tradeUserEvent) Add() *tradeUserEvent {
	base.ParentNames = []string{base.EventName, "ProductEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateTradeUserRegisteredEvent(bullionId string, tradeUser *interfaces.TradeUserEntity, adminId string) *BaseEvent {
	event := &tradeUserEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUser.ID,
			AdminId:   adminId,
			Payload:   tradeUser,
			EventName: "TradeUserRegisteredEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateTradeUserActivatedEvent(bullionId string, tradeUser *interfaces.TradeUserEntity, adminId string) *BaseEvent {
	event := &tradeUserEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUser.ID,
			AdminId:   adminId,
			Payload:   tradeUser,
			EventName: "TradeUserActivatedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateTradeUserDisabledEvent(bullionId string, tradeUser *interfaces.TradeUserEntity, adminId string) *BaseEvent {
	event := &tradeUserEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUser.ID,
			AdminId:   adminId,
			Payload:   tradeUser,
			EventName: "TradeUserDisabledEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateTradeUserMarginModifiedEvent(bullionId string, tradeUser *interfaces.TradeUserEntity, adminId string) *BaseEvent {
	event := &tradeUserEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUser.ID,
			AdminId:   adminId,
			Payload:   tradeUser.AvailableMargins,
			EventName: "TradeUserMarginModifiedEvent",
		},
	}
	event.Add()
	return event.BaseEvent
}
