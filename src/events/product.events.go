package events

import "github.com/rpsoftech/bullion-server/src/interfaces"

type productEvent struct {
	*BaseEvent
}

func (base *productEvent) Add() *productEvent {
	base.ParentNames = append(base.ParentNames, "ProductEvent")
	base.BaseEvent.CreateBaseEvent()
	return base
}

type productCreatedEvent struct {
	*productEvent
}

func CreateProductCreatedEvent(bullionId string, productId string, product *interfaces.ProductEntity, adminId string) *productCreatedEvent {
	event := &productCreatedEvent{
		productEvent: &productEvent{
			BaseEvent: &BaseEvent{
				BullionId:   bullionId,
				KeyId:       productId,
				AdminId:     adminId,
				Payload:     product,
				EventName:   "ProductCreatedEvent",
				ParentNames: []string{"ProductCreatedEvent"},
			},
		},
	}
	event.Add()
	return event
}
