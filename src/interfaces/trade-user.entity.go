package interfaces

type (
	TradeUserBase struct {
		BullionId   string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Name        string `bson:"name" json:"name" validate:"required,min=2,max=100"`
		Number      string `bson:"number" json:"number" validate:"required,min=12,max=12"`
		Email       string `bson:"email" json:"email" validate:"required"`
		CompanyName string `bson:"companyName" json:"companyName" validate:"required,min=2,max=50"`
		GstNumber   string `bson:"gstNumber" json:"gstNumber" validate:"required,gstNumber"`
		RawPassword string `bson:"rawPassword" json:"password" mapstructure:"password" validate:"required,min=4"`
	}

	TradeUserAdvanced struct {
		UserName string `bson:"userName" json:"userName" validate:"required"`
		IsActive bool   `bson:"isActive" json:"isActive" validate:"boolean"`
		UNumber  string `bson:"uNumber" json:"uNumber" validate:"required"`
	}
	UserMarginsDataStruct struct {
		Gold   int32 `bson:"gold" json:"gold" validate:"min=0"`
		Silver int32 `bson:"silver" json:"silver" validate:"min=0"`
	}

	TradeUserMargins struct {
		AllotedMargins   *UserMarginsDataStruct `bson:"allotedMargins" json:"allotedMargins" validate:"required"`
		AvailableMargins *UserMarginsDataStruct `bson:"availableMargins" json:"availableMargins" validate:"required"`
	}

	TradeUserEntity struct {
		*BaseEntity        `bson:"inline"`
		*TradeUserBase     `bson:"inline"`
		*passwordEntity    `bson:"inline"`
		*TradeUserAdvanced `bson:"inline"`
		*TradeUserMargins  `bson:"margins" json:"margins"`
	}

	ApiTradeUserRegisterResponse struct {
		UserToken   string `json:"userToken"`
		OtpReqToken string `json:"otpReqToken"`
	}
)

func (user *TradeUserEntity) CreateNew() (r *TradeUserEntity) {
	user.createNewId()
	return user
}

func (user *TradeUserEntity) UpdateUser() (r *TradeUserEntity) {
	user.passwordEntity = CreatePasswordEntity(user.RawPassword)
	user.Updated()
	return user
}
