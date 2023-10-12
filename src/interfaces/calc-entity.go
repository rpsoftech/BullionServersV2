package interfaces

type CshPremiumBuySellSnapshot struct {
	Tax     int     `bson:"tax" json:"tax" validate:"required"`
	Premium float32 `bson:"premium" json:"premium" validate:"required"`
	// Tcs     float32 `bson:"tcs" json:"tcs" validate:"required"`
	// Tds     float32 `bson:"tds" json:"tds" validate:"required"`
}

type CalcSnapshotStruct struct {
	Buy  CshPremiumBuySellSnapshot `bson:"buy" json:"buy" validate:"required"`
	Sell CshPremiumBuySellSnapshot `bson:"sell" json:"sell" validate:"required"`
}

type CalcEntity struct {
	*BaseEntity         `bson:"inline"`
	ProductId           string `bson:"productId" json:"productId" validate:"required"`
	*CalcSnapshotStruct `bson:"inline"`
}
