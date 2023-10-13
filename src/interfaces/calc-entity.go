package interfaces

type CshPremiumBuySellSnapshot struct {
	Tax     int     `bson:"tax" json:"tax" validate:"min=0,max=50"`
	Premium float32 `bson:"premium" json:"premium"`
}

type CalcSnapshotStruct struct {
	Buy  CshPremiumBuySellSnapshot `bson:"buy" json:"buy" validate:"required"`
	Sell CshPremiumBuySellSnapshot `bson:"sell" json:"sell" validate:"required"`
}

// type CalcEntity struct {
// 	*BaseEntity         `bson:"inline"`
// 	ProductId           string `bson:"productId" json:"productId" validate:"required"`
// 	*CalcSnapshotStruct `bson:"inline"`
// }
