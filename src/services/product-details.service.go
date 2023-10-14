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
		productRepo:   repos.Pro
		

func (service *productService) AddNewProduct(productBase *interfaces.ProductBaseStruct, calcBase *interfaces.CalcSnapshotStruct) error {
	currentProducts, err := service.GetProductsByBullionId(productBase.BullionId)
	if err != nil {
		return err
	}
	service.saveProductEntity(entity)
	return nil
}

func (service *productService) saveProductEntity(entity *interfaces.ProductEntity) {
	service.productRepo.Save(entity)
	
}

func (service *productService) saveProductEntityToLocalCaches(entity *interfaces.ProductEntity) {
	if _, ok := service.productsById[entity.BullionId]; !ok {
		service.productsById[entity.BullionId] = make(map[string]*interfaces.ProductEntity)
	}
	service.productsById[entity.BullionId][entity.ID] = entity
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
	return products, nil
}