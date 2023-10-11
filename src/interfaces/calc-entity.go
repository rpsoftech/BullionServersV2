package interfaces

type CshPremiumBuySellSnapshot struct {
	Tax     float32 `bson:"tax" json:"tax" validate:"required"`
	Tcs     float32 `bson:"tcs" json:"tcs" validate:"required"`
	Tds     float32 `bson:"tds" json:"tds" validate:"required"`
	Premium float32 `bson:"premium" json:"premium" validate:"required"`
}

type CalcSnapshotStruct struct {
	Buy  CshPremiumBuySellSnapshot `bson:"buy" json:"buy" validate:"required"`
	Sell CshPremiumBuySellSnapshot `bson:"sell" json:"sell" validate:"required"`
}

type CalcEntity struct {
	*BaseEntity         `bson:"inline"`
	*CalcSnapshotStruct `bson:"inline"`
}
