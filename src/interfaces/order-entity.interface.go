package interfaces

type (
	OrderEntity struct {
		*BaseEntity           `bson:"inline"`
		*OrderBase            `bson:"inline"`
		*LimitWatcherRequired `bson:"inline"`
		*FromAdmin            `bson:"inline,omitempty"`
	}

	LimitWatcherRequired struct {
		Price             float64 `json:"price" bson:"price" validate:"required"`
		ProductId         string  `json:"productId" bson:"productId" validate:"required,uuid"`
		GroupId           string  `json:"groupId" bson:"groupId" validate:"required,uuid"`
		ProductGroupMapId string  `json:"productGroupMapId" bson:"productGroupMapId" validate:"required,uuid"`
	}

	FromAdmin struct {
		PlacedBy string `bson:"placedBy,omitempty" json:"placedBy,omitempty" validate:"uuid"`
		AutoRate bool   `bson:"autoRate,omitempty" json:"autoRate,omitempty" validate:"boolean"`
	}

	OrderBase struct {
		UserId      string      `bson:"userId" json:"userId" validate:"required,uuid"`
		OrderType   OrderType   `bson:"orderType" json:"orderType" validate:"required,enum=OrderStatus"`
		OrderStatus OrderStatus `bson:"orderStatus" json:"orderStatus" validate:"required,enum=OrderStatus"`
		BuySell     BuySell     `bson:"buySell" json:"buySell" validate:"required,enum=BuySell"`
	}
)
