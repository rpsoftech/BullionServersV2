package interfaces

type (
	FeedsBase struct {
		BullionId string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Title     string `bson:"title" json:"title" validate:"required,min=3"`
		Body      string `bson:"body" json:"body" validate:"required,min=3"`
		IsHtml    bool   `bson:"isHtml" json:"isHtml" validate:"boolean"`
	}
	FeedsEntity struct {
		*BaseEntity `bson:"inline"`
		*FeedsBase  `bson:"inline"`
	}
	FeedUpdateRequestBody struct {
		FeedId     string `bson:"feedId" json:"feedId" validate:"required,uuid"`
		*FeedsBase `bson:"inline"`
	}
)

func (b *FeedsEntity) CreateNewId() *FeedsEntity {
	if b.BaseEntity == nil {
		b.BaseEntity = &BaseEntity{}
	}
	b.createNewId()
	return b
}
