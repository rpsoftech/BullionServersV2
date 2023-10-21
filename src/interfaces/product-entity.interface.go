package interfaces

type ProductBaseStruct struct {
	BullionId           string               `bson:"bullionId" json:"bullionId" validate:"required"`
	Name                string               `bson:"name" json:"name" validate:"required"`
	SourceSymbol        SymbolsEnum          `bson:"sourceSymbol" json:"sourceSymbol" validate:"required,enum=SymbolsEnum"`
	CalculationSymbol   SymbolsEnum          `bson:"calculationSymbol" json:"calculationSymbol" validate:"required,enum=SymbolsEnum"`
	IsActive            bool                 `bson:"isActive" json:"isActive" validate:"boolean"`
	IsHedging           bool                 `bson:"isHedging" json:"isHedging" validate:"boolean"`
	FloatPoint          int                  `bson:"floatPoint" json:"floatPoint" validate:"min=0,max=4"`
	CalculatedOnPriceOf CalculateOnPriceType `bson:"calculatedOnPriceOf" json:"calculatedOnPriceOf" validate:"required,enum=CalculateOnPriceType"`
}

type ProductEntity struct {
	*BaseEntity        `bson:"inline"`
	*ProductBaseStruct `bson:"inline"`
	Sequence           int                 `bson:"sequence" json:"sequence"`
	CalcSnapshot       *CalcSnapshotStruct `bson:"calcSnapshot" json:"calcSnapshot" validate:"required"`
}

func CreateNewProduct(productBase *ProductBaseStruct, calcSnapShot *CalcSnapshotStruct, sequence int) (r *ProductEntity) {
	b := &ProductEntity{
		ProductBaseStruct: productBase,
		CalcSnapshot:      calcSnapShot,
		Sequence:          sequence,
		BaseEntity:        &BaseEntity{},
	}
	b.createNewId()
	return b
}

type UpdateProductApiBody struct {
	ProductId          string `json:"id" validate:"required,uuid"`
	*ProductBaseStruct `validate:"required"`
	CalcSnapshot       *CalcSnapshotStruct `json:"calcSnapShot" validate:"required"`
}
