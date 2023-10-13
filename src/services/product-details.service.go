package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type productService struct {
	productRepo *repos.ProductRepoStruct
	products    map[string]map[string]*interfaces.ProductEntity
}

var ProductService *productService

func init() {
	ProductService = &productService{
		productRepo: repos.ProductRepo,
		products:    make(map[string]map[string]*interfaces.ProductEntity),
	}
}

func (service *productService) AddNewProduct(prouctBase *interfaces.ProductBaseStruct, calcBase *interfaces.CalcSnapshotStruct, bullionId string) {
	// if

}
func (service *productService) GetProductsByBullionId(bullionId string) {

}
