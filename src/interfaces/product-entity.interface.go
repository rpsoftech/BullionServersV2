package interfaces

type ProductBaseStruct struct {
	BullionId           string               `bson:"bullionId" json:"bullionId" validate:"required"`
	Name                string               `bson:"name" json:"name" validate:"required"`
	SourceSymbol        string               `bson:"sourceSymbol" json:"sourceSymbol" validate:"required"`
	CalculationSymbol   string               `bson:"CalculationSymbol" json:"CalculationSymbol" validate:"required"`
	IsActive            bool                 `bson:"isActive" json:"isActive" validate:"required"`
	IsHedging           bool                 `bson:"isHedging" json:"isHedging" validate:"required"`
	FloatPoint          int                  `bson:"floatPoint" json:"floatPoint" validate:"required"`
	CalculatedOnPriceOf CalculateOnPriceType `bson:"calculatedOnPriceOf" json:"calculatedOnPriceOf" validate:"required"`
}

type ProductEntity struct {
	*BaseEntity        `bson:"inline"`
	*ProductBaseStruct `bson:"inline"`
	CalcSnapshot       *CalcSnapshotStruct `bson:"calcSnapshot" json:"calcSnapshot" validate:"required"`
}

func CreateNewProduct(productBase *ProductBaseStruct, calcSnapShot *CalcSnapshotStruct) (r *ProductEntity) {
	b := &ProductEntity{
		ProductBaseStruct: productBase,
		CalcSnapshot:      calcSnapShot,
		BaseEntity:        &BaseEntity{},
	}
	b.createNewId()
	return b
}
