package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type FlagServiceStruct struct {
	firebaseDb *firebaseDatabaseService
	eventBus   *eventBusService
	flagsMap   map[string]*interfaces.FlagsInterface
}

var FlagService *FlagServiceStruct

func init() {
	getFlagService()
}

func getFlagService() *FlagServiceStruct {
	if FlagService == nil {
		FlagService = &FlagServiceStruct{
			firebaseDb: getFirebaseRealTimeDatabase(),
			eventBus:   getEventBusService(),
			flagsMap:   make(map[string]*interfaces.FlagsInterface),
		}
		println("Flag Service Initialized")
	}
	return FlagService
}

func (s *FlagServiceStruct) UpdateFlags(entity *interfaces.FlagsInterface, adminId string) (*interfaces.FlagsInterface, error) {
	if err := s.firebaseDb.SetPublicData(entity.BullionId, []string{"flags"}, entity); err != nil {
		return nil, err
	}
	s.eventBus.Publish(events.FlagsUpdatedEvent(entity, adminId))
	s.flagsMap[entity.BullionId] = entity
	return entity, nil
}

func (s *FlagServiceStruct) GetFlags(bullionId string) (*interfaces.FlagsInterface, error) {
	if entity := s.flagsMap[bullionId]; entity != nil {
		return entity, nil
	}

	entity := new(interfaces.FlagsInterface)
	if err := s.firebaseDb.GetPublicData(bullionId, []string{"flags"}, entity); err != nil {
		return nil, err
	}
	s.flagsMap[bullionId] = entity
	return entity, nil
}
