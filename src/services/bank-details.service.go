package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type bankDetailsService struct {
	bankDetailsRepo *repos.BankDetailsRepoStruct
	eventBus        *eventBusService
}

var BankDetailsService *bankDetailsService

func init() {
	BankDetailsService = &bankDetailsService{
		eventBus:        EventBus,
		bankDetailsRepo: repos.BankDetailsRepo,
	}
	if EventBus == nil {
		go func() {
			if EventBus == nil {
				panic("Event Bus Not Initialized")
			}
			BankDetailsService.eventBus = EventBus
		}()
	}
	println("Bank Details Service Initialized")
}

func (s *bankDetailsService) GetBankDetailsByBullionId(id string) (*[]interfaces.BankDetailsEntity, error) {
	return s.bankDetailsRepo.GetAllByBullionId(id)
}

func (s *bankDetailsService) addUpdateBankDetails(entity *interfaces.BankDetailsEntity) (*interfaces.BankDetailsEntity, error) {
	_, err := s.bankDetailsRepo.Save(entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}
func (s *bankDetailsService) UpdateBankDetails(entity *interfaces.BankDetailsBase, id string, adminId string) (*interfaces.BankDetailsEntity, error) {
	entityFromDb, err := s.bankDetailsRepo.FindOne(id)
	if err != nil {
		return nil, err
	}
	if entityFromDb.BullionId != entity.BullionId {
		return nil, &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "You do not have access to this Bank Details",
			Name:       "ERROR_MISMATCH_BULLION_ID",
		}
	}
	entityFromDb.BankDetailsBase = entity
	_, err = s.addUpdateBankDetails(entityFromDb)
	if err != nil {
		return nil, err
	}
	event := events.CreateBankDetailsUpdatedEvent(entityFromDb, adminId)
	s.eventBus.Publish(event)
	return entityFromDb, err
}
func (s *bankDetailsService) AddNewBankDetails(base *interfaces.BankDetailsBase, adminId string) (*interfaces.BankDetailsEntity, error) {
	entity := interfaces.CreateNewBankDetails(base)
	_, err := s.addUpdateBankDetails(entity)
	if err != nil {
		return nil, err
	}
	event := events.CreateNewBankDetailsCreated(entity, adminId)
	s.eventBus.Publish(event)
	return entity, err
}
func (s *bankDetailsService) DeleteBankDetails(entity *interfaces.BankDetailsBase, id string, adminId string) (*interfaces.BankDetailsEntity, error) {
	entityFromDb, err := s.bankDetailsRepo.FindOne(id)
	if err != nil {
		return nil, err
	}
	if entityFromDb.BullionId != entity.BullionId {
		return nil, &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "You do not have access to this Bank Details",
			Name:       "ERROR_MISMATCH_BULLION_ID",
		}
	}
	err = s.bankDetailsRepo.DeleteById(id)
	if err != nil {
		return nil, err
	}
	event := events.CreateBankDetailsDeletedEvent(entity, id, adminId)
	s.eventBus.Publish(event)
	return entityFromDb, err
}
