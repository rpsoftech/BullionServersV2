package interfaces

type (
	MsgTemplateBase struct {
		BullionId        string `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		WhatsappTemplate string `bson:"whatsappTemplate" json:"whatsappTemplate" validate:"required"`
		MSG91Id          string `bson:"msg91Id" json:"msg91Id" validate:"required"`
	}
)
