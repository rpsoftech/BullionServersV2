package interfaces

type (
	TradeUserBase struct {
		BullionId   string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Name        string `bson:"name" json:"name" validate:"required,min=2,max=100"`
		Number      string `bson:"number" json:"number" validate:"required,min=10,max=10"`
		Email       string `bson:"email" json:"email" validate:"required"`
		CompanyName string `bson:"companyName" json:"companyName" validate:"required,min=2,max=50"`
		GstNumber   string `bson:"gstNumber" json:"gstNumber" validate:"required,gstNumber"`
	}

	TradeUserAdvanced struct {
		UserName string `bson:"userName" json:"userName" validate:"required"`
		IsActive bool   `bson:"isActive" json:"isActive" validate:"required"`
		UNumber  string `bson:"uNumber" json:"uNumber" validate:"required"`
	}
	UserMarginsDataStruct struct {
		Gold   int32 `bson:"gold" json:"gold" validate:"min=0"`
		Silver int32 `bson:"silver" json:"silver" validate:"min=0"`
	}

	TradeUserMargins struct {
		AllotedMargins   UserMarginsDataStruct `bson:"allotedMargins" json:"allotedMargins" validate:"required"`
		AvailableMargins UserMarginsDataStruct `bson:"availableMargins" json:"availableMargins" validate:"required"`
	}

	TradeUserEntity struct {
		*BaseEntity        `bson:"inline"`
		*TradeUserBase     `bson:"inline"`
		*passwordEntity    `bson:"inline"`
		*TradeUserAdvanced `bson:"inline"`
		*TradeUserMargins  `bson:"inline"`
	}
)

func (user *TradeUserEntity) CreateNew() (r *TradeUserEntity) {
	user.createNewId()
	return user
}
