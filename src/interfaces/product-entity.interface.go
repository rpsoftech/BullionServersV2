package interfaces

type ProductEntity struct {
	*BaseEntity         `bson:"inline"`
	BullionId           string               `bson:"bullionId" json:"bullionId" validate:"required"`
	Name                string               `bson:"name" json:"name" validate:"required"`
	SourceSymbol        string               `bson:"sourceSymbol" json:"sourceSymbol" validate:"required"`
	CalculationSymbol   string               `bson:"CalculationSymbol" json:"CalculationSymbol" validate:"required"`
	IsActive            bool                 `bson:"isActive" json:"isActive" validate:"required"`
	IsHedging           bool                 `bson:"isHedging" json:"isHedging" validate:"required"`
	CalculatedOnPriceOf CalculateOnPriceType `bson:"calculatedOnPriceOf" json:"calculatedOnPriceOf" validate:"required"`
	CalcSnapshot        CalcSnapshotStruct   `bson:"calcSnapshot" json:"calcSnapshot" validate:"required"`
}
