package interfaces

type (
	TradeUserBase struct {
		Name        string `bson:"name" json:"name" validate:"required,min=2,max=100"`
		Number      string `bson:"number" json:"number" validate:"required,min=10,max=10"`
		CompanyName string `bson:"companyName" json:"companyName" validate:"required,min=2,max=30"`
		GstNumber   string `bson:"gstNumber" json:"gstNumber" validate:"required,gstNumber"`
		// Password        string `bson:"password" json:"password" validate:"required"`
		*passwordEntity `bson:"inline"`
	}

	TradeUserAdvanced struct {
		UserName string `bson:"userName" json:"userName" validate:"required"`
		IsActive bool   `bson:"isActive" json:"isActive" validate:"required"`
		UNo      string `bson:"uNo" json:"uNo" validate:"required"`
	}

	TradeUserMargins struct {
		AllotedMargins   UserMarginsDataStruct `bson:"allotedMargins" json:"allotedMargins" validate:"required"`
		AvailableMargins UserMarginsDataStruct `bson:"availableMargins" json:"availableMargins" validate:"required"`
	}

	UserMarginsDataStruct struct {
		Gold   int32 `bson:"gold" json:"gold" validate:"min=0"`
		Silver int32 `bson:"silver" json:"silver" validate:"min=0"`
	}

	TradeUserEntity struct {
		*BaseEntity        `bson:"inline"`
		*TradeUserBase     `bson:"inline"`
		*TradeUserAdvanced `bson:"inline"`
		*TradeUserMargins  `bson:"inline"`
	}
)
