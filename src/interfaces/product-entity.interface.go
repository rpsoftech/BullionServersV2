package interfaces

// id: ProductID;
// name: string;
// sourceSymbole: SourceSymbole;
// calculationSymbole: CaculationSymbole;
// isActive: boolean;
// isHedging: boolean;
// showLocation: ProductShowLocation;
// calculatedOnPriceof: CalculatedOnPriceof;
// calculatedOnPriceType: CalculatedOnPriceType;
// calcSnapshot: CshVariableSnapshot;
// createdAt: Date;
// modifiedAt: Date;

type ProductEntity struct {
	*BaseEntity `bson:"inline"`
	Name        string `bson:"name" json:"name" validate:"required"`
	BaseSymbol  string `bson:"baseSymbol" json:"baseSymbol" validate:"required"`
}
