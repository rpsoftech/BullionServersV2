package interfaces

type (
	TradeUserGroupBase struct {
		BullionId string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Name      string `bson:"name" json:"name" validate:"required"`
		IsActive  bool   `bson:"isActive" json:"isActive" validate:"boolean"`
		CanTrade  bool   `bson:"canTrade" json:"canTrade" validate:"boolean"`
		CanLogin  bool   `bson:"canLogin" json:"canLogin" validate:"boolean"`
	}
	TradeUserGroupEntity struct {
		*BaseEntity         `bson:"inline"`
		*TradeUserGroupBase `bson:"inline"`
	}

	TradeUserGroupMapEntity struct {
		*BaseEntity            `bson:"inline"`
		*TradeUserGroupMapBase `bson:"inline"`
	}

	TradeUserGroupMapBase struct {
		BullionId         string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		GroupId           string `bson:"groupId" json:"groupId" validate:"required,uuid"`
		ProductId         string `bson:"productId" json:"productId" validate:"required,uuid"`
		IsActive          bool   `bson:"isActive" json:"isActive" validate:"boolean"`
		CanTrade          bool   `bson:"canTrade" json:"canTrade" validate:"boolean"`
		*GroupPremiumBase `bson:"groupPremiumBase" json:"groupPremiumBase" validate:"required"`
		*GroupVolumeBase  `bson:"groupVolumeBase" json:"groupVolumeBase" validate:"required"`
	}

	GroupPremiumBase struct {
		Buy  float64 `bson:"buy" json:"buy"`
		Sell float64 `bson:"sell" json:"sell"`
	}

	GroupVolumeBase struct {
		OneClick int `bson:"oneClick" json:"oneClick"`
		Step     int `bson:"step" json:"step"`
		Total    int `bson:"total" json:"total"`
	}
)

func (r *TradeUserGroupEntity) CreateNew() *TradeUserGroupEntity {
	r.BaseEntity = &BaseEntity{}
	r.createNewId()
	return r
}

func (r *TradeUserGroupMapEntity) CreateNew() *TradeUserGroupMapEntity {
	r.BaseEntity = &BaseEntity{}
	r.createNewId()
	return r
}
