package services

import (
	"fmt"
	"reflect"

	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type (
	productService struct {
		productRepo                   *repos.ProductRepoStruct
		firebaseDatabaseService       *firebaseDatabaseService
		eventBus                      *eventBusService
		productsByBullionAndProductId map[string]map[string]*interfaces.ProductEntity
		productsArray                 map[string]*[]interfaces.ProductEntity
		productsById                  map[string]*interfaces.ProductEntity
	}
)

var ProductService *productService

func init() {
	ProductService = &productService{
		eventBus:                      getEventBusService(),
		firebaseDatabaseService:       getFirebaseRealTimeDatabase(),
		productRepo:                   repos.ProductRepo,
		productsByBullionAndProductId: make(map[string]map[string]*interfaces.ProductEntity),
		productsArray:                 make(map[string]*[]interfaces.ProductEntity),
		productsById:                  make(map[string]*interfaces.ProductEntity),
	}
	println("Product Service Initialized")
}

func (service *productService) AddNewProduct(productBase *interfaces.ProductBaseStruct, calcBase *interfaces.CalcSnapshotStruct, adminId string) (*interfaces.ProductEntity, error) {
	currentProducts, err := service.GetProductsByBullionId(productBase.BullionId)
	if err != nil {
		return nil, err
	}
	currentCount := len(*currentProducts)
	entity := interfaces.CreateNewProduct(productBase, calcBase, currentCount+1)
	// TODO: Create Product Group Map
	_, err = service.saveProductEntity(entity)
	if err != nil {
		return nil, err
	}
	event := events.CreateProductCreatedEvent(entity.BullionId, entity.ID, entity, adminId)
	service.eventBus.Publish(event.BaseEvent)
	return entity, nil
}

func (service *productService) UpdateProductSequence(updateProductCalcSequenceApiBody *[]interfaces.UpdateProductCalcSequenceApiBody, adminId string, bullionID string) (*[]interfaces.ProductEntity, error) {
	entities := make([]interfaces.ProductEntity, len(*updateProductCalcSequenceApiBody))
	modified := make([]bool, len(*updateProductCalcSequenceApiBody))
	for i, prod := range *updateProductCalcSequenceApiBody {
		oldDetails, err := service.GetProductsById(prod.BullionId, prod.ProductId)
		if err != nil {
			return nil, err
		}
		if oldDetails.Sequence != prod.Sequence {
			modified[i] = true
		} else {
			modified[i] = false
		}
		oldDetails.Sequence = prod.Sequence
		entities[i] = *oldDetails
	}
	result, err := service.productRepo.BulkUpdate(&entities)
	if err == nil {
		event := events.CreateProductSequenceChangedEvent(bullionID, &entities, adminId)
		service.eventBus.PublishAll(event)
		for i, entity := range entities {
			if modified[i] {
				service.saveProductEntityToLocalCaches(&entity, true)
			}
		}
	}
	return result, err
}

func (service *productService) UpdateProductCalcSnapshot(updateProductCalcSnapshot *[]interfaces.UpdateProductCalcSnapshotApiBody, adminId string) (*[]interfaces.ProductEntity, error) {
	entities := make([]interfaces.ProductEntity, len(*updateProductCalcSnapshot))
	modified := make([]bool, len(*updateProductCalcSnapshot))
	for i, prod := range *updateProductCalcSnapshot {
		oldDetails, err := service.GetProductsById(prod.BullionId, prod.ProductId)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(oldDetails.CalcSnapshot, prod.CalcSnapshot) {
			modified[i] = true
		} else {
			modified[i] = false
		}
		oldDetails.CalcSnapshot = prod.CalcSnapshot
		entities[i] = *oldDetails
	}
	result, err := service.productRepo.BulkUpdate(&entities)
	if err == nil {
		for i, entity := range entities {
			if modified[i] {
				event := events.CreateProductCalcUpdated(entity.BullionId, entity.ID, entity.CalcSnapshot, adminId)
				service.saveProductEntityToLocalCaches(&entity, true)
				service.eventBus.Publish(event.BaseEvent)
			}
		}
	}
	return result, err
}

func (service *productService) UpdateProduct(updateProductBody *[]interfaces.UpdateProductApiBody, adminId string) (*[]interfaces.ProductEntity, error) {
	entities := make([]interfaces.ProductEntity, len(*updateProductBody))
	modified := make([]bool, len(*updateProductBody))
	for i, prod := range *updateProductBody {
		oldDetails, err := service.GetProductsById(prod.BullionId, prod.ProductId)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(oldDetails.ProductBaseStruct, prod.ProductBaseStruct) || !reflect.DeepEqual(oldDetails.CalcSnapshot, prod.CalcSnapshot) {
			modified[i] = true
		} else {
			modified[i] = false
		}
		oldDetails.ProductBaseStruct = prod.ProductBaseStruct
		oldDetails.CalcSnapshot = prod.CalcSnapshot
		entities[i] = *oldDetails
	}
	result, err := service.productRepo.BulkUpdate(&entities)
	if err == nil {
		for i, entity := range entities {
			if modified[i] {
				service.saveProductEntityToLocalCaches(&entity, true)
				event := events.CreateProductUpdatedEvent(entity.BullionId, entity.ID, &entity, adminId)
				service.eventBus.Publish(event.BaseEvent)
			}
		}
	}
	return result, err
}

func (service *productService) saveProductEntity(entity *interfaces.ProductEntity) (*interfaces.ProductEntity, error) {
	_, err := service.productRepo.Save(entity)
	if err != nil {
		return entity, err
	}
	service.saveProductEntityToLocalCaches(entity, true)
	return entity, nil
}

func (service *productService) saveProductEntityToLocalCaches(entity *interfaces.ProductEntity, appendToArray bool) {
	service.firebaseDatabaseService.SetPublicData(entity.BullionId, []string{"products", entity.ID}, entity)
	if _, ok := service.productsByBullionAndProductId[entity.BullionId]; !ok {
		service.productsByBullionAndProductId[entity.BullionId] = make(map[string]*interfaces.ProductEntity)
	}
	service.productsByBullionAndProductId[entity.BullionId][entity.ID] = entity
	service.productsById[entity.ID] = entity

	if !appendToArray {
		return
	}

	if _, ok := service.productsArray[entity.BullionId]; !ok {
		service.productsArray[entity.BullionId] = &[]interfaces.ProductEntity{}
	} else {
		found := false
		for index, ele := range *service.productsArray[entity.BullionId] {
			if ele.ID == entity.ID {
				found = true
				(*service.productsArray[entity.BullionId])[index] = *entity
			}
		}
		if !found {
			*service.productsArray[entity.BullionId] = append(*service.productsArray[entity.BullionId], *entity)
		}
	}
}

func (service *productService) GetProductsByBullionId(bullionId string) (*[]interfaces.ProductEntity, error) {
	if result, ok := service.productsArray[bullionId]; ok {
		return result, nil
	}
	products, err := service.productRepo.FindByBullionId(bullionId)
	if err != nil {
		return nil, err
	}
	service.productsArray[bullionId] = products
	for _, product := range *products {
		service.saveProductEntityToLocalCaches(&product, false)
	}
	return products, nil
}

func (service *productService) GetProductsById(bullionId string, productId string) (*interfaces.ProductEntity, error) {
	allProducts, err := service.GetProductsByBullionId(bullionId)
	if err != nil {
		return nil, err
	}
	for _, product := range *allProducts {
		if product.ID == productId {
			return &product, nil
		}
	}
	return nil, &interfaces.RequestError{
		StatusCode: 400,
		Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
		Message:    fmt.Sprintf("Product Entities identified by bullionId %s and productId %s not found", bullionId, productId),
		Name:       "ENTITY_NOT_FOUND",
	}
}
