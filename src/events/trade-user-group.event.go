package events

import "github.com/rpsoftech/bullion-server/src/interfaces"

type tradeUserGroupEvent struct {
	*BaseEvent `bson:"inline"`
}

func (base *tradeUserGroupEvent) Add() *tradeUserGroupEvent {
	base.ParentNames = []string{base.EventName, "ProductEvent"}
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateTradeUserGroupCreated(bullionId string, tradeUserGroup *interfaces.TradeUserGroupEntity, adminId string) *BaseEvent {
	event := &tradeUserGroupEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUserGroup.ID,
			AdminId:   adminId,
			Payload:   tradeUserGroup,
			EventName: "TradeUserGroupCreated",
		},
	}
	event.Add()
	return event.BaseEvent
}
func CreateTradeUserGroupUpdated(bullionId string, tradeUserGroup *interfaces.TradeUserGroupEntity, adminId string) *BaseEvent {
	event := &tradeUserGroupEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     tradeUserGroup.ID,
			AdminId:   adminId,
			Payload:   tradeUserGroup,
			EventName: "TradeUserGroupUpdated",
		},
	}
	event.Add()
	return event.BaseEvent
}

func CreateTradeUserGroupMapUpdated(bullionId string, tradeUserGroupMap *[]interfaces.TradeUserGroupMapEntity, groupId string, adminId string) *BaseEvent {
	event := &tradeUserGroupEvent{
		BaseEvent: &BaseEvent{
			BullionId: bullionId,
			KeyId:     groupId,
			AdminId:   adminId,
			Payload:   tradeUserGroupMap,
			EventName: "TradeUserGroupMapUpdated",
		},
	}
	event.Add()
	return event.BaseEvent
}
