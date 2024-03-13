package events

import "github.com/rpsoftech/bullion-server/src/interfaces"

type productEvent struct {
	*BaseEvent `bson:"inline"`
}

type productSequenceChangedEvent struct {
	Id        string `bson:"id" json:"id"`
	BullionId string `bson:"bullionId" json:"bullionId"`
	Sequence  int    `bson:"sequence" json:"sequence"`
}

func (base *productEvent) Add() *productEvent {
	base.ParentNames = append(base.ParentNames, "ProductEvent")
	base.BaseEvent.CreateBaseEvent()
	return base
}

func CreateProductCreatedEvent(bullionId string, productId string, product *interfaces.ProductEntity, adminId string) *productEvent {
	event := &productEvent{
		BaseEvent: &BaseEvent{
			BullionId:   bullionId,
			KeyId:       productId,
			AdminId:     adminId,
			Payload:     product,
			EventName:   "ProductCreatedEvent",
			ParentNames: []string{"ProductCreatedEvent"},
		},
	}
	event.Add()
	return event
}

func CreateProductUpdatedEvent(bullionId string, productId string, product *interfaces.ProductEntity, adminId string) *productEvent {
	event := &productEvent{
		BaseEvent: &BaseEvent{
			BullionId:   bullionId,
			KeyId:       productId,
			AdminId:     adminId,
			Payload:     product,
			EventName:   "ProductUpdatedEvent",
			ParentNames: []string{"ProductUpdatedEvent"},
		},
	}
	event.Add()
	return event
}

func CreateProductSequenceChangedEvent(bullionId string, product *[]interfaces.ProductEntity, adminId string) *[]BaseEvent {
	events := make([]BaseEvent, len(*product))
	for i, pro := range *product {
		event := productEvent{
			BaseEvent: &BaseEvent{
				BullionId: bullionId,
				KeyId:     pro.ID,
				AdminId:   adminId,
				Payload: productSequenceChangedEvent{
					Id:        pro.ID,
					BullionId: pro.BullionId,
					Sequence:  pro.Sequence,
				},
				EventName:   "ProductSequenceChangedEvent",
				ParentNames: []string{"ProductSequenceChangedEvent"},
			},
		}
		event.Add()
		events[i] = *event.BaseEvent
	}
	return &events
}

func CreateProductCalcUpdated(bullionId string, productId string, calcSnapshot *interfaces.CalcSnapshotStruct, adminId string) *productEvent {
	event := &productEvent{
		BaseEvent: &BaseEvent{
			BullionId:   bullionId,
			KeyId:       productId,
			AdminId:     adminId,
			Payload:     calcSnapshot,
			EventName:   "ProductCalcUpdated",
			ParentNames: []string{"ProductCalcUpdated"},
		},
	}
	event.Add()
	return event
}

func CreateProductDisabled(bullionId string, productId string, adminId string) *productEvent {
	event := &productEvent{
		BaseEvent: &BaseEvent{
			BullionId:   bullionId,
			KeyId:       productId,
			AdminId:     adminId,
			Payload:     "",
			EventName:   "ProductDisabled",
			ParentNames: []string{"ProductDisabled"},
		},
	}
	event.Add()
	return event
}
