package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type productService struct {
	productRepo   *repos.ProductRepoStruct
	productsById  map[string]map[string]*interfaces.ProductEntity
	productsArray map[string]*[]interfaces.ProductEntity
}

var ProductService *productService

func init() {
	ProductService = &productService{
		productRepo:   repos.ProductRepo,
		productsById:  make(map[string]map[string]*interfaces.ProductEntity),
		productsArray: make(map[string]*[]interfaces.ProductEntity),
	}
}

func (service *productService) AddNewProduct(productBase *interfaces.ProductBaseStruct, calcBase *interfaces.CalcSnapshotStruct) (*interfaces.ProductEntity, error) {
	currentProducts, err := service.GetProductsByBullionId(productBase.BullionId)
	if err != nil {
		return nil, err
	}
	currentCount := len(*currentProducts)
	entity := interfaces.CreateNewProduct(productBase, calcBase)
	entity.Sequence = currentCount + 1
	return service.saveProductEntity(entity)
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
	if _, ok := service.productsById[entity.BullionId]; !ok {
		service.productsById[entity.BullionId] = make(map[string]*interfaces.ProductEntity)
	}
	service.productsById[entity.BullionId][entity.ID] = entity

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
