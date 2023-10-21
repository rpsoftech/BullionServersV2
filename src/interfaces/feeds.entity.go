package interfaces

type (
	FeedsEntity struct {
		*BaseEntity `bson:"inline"`
		BullionId   string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Title       string `bson:"title" json:"title" validate:"required,min=3"`
		Body        string `bson:"body" json:"body" validate:"required,min=3"`
		IsHtml      bool   `bson:"isHtml" json:"isHtml" validate:"boolean"`
	}
)
