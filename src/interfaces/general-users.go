package interfaces

type GeneralUser struct {
	BaseEntity         `bson:"inline"`
	UserRolesInterface `bson:"inline"`
	FirstName          string     `bson:"firstName" json:"firstName" validate:"required"`
	LastName           string     `bson:"lastName" json:"lastName" validate:"required"`
	RandomPass         string     `bson:"randomPass" json:"randomPass" validate:"required"`
	FirmName           string     `bson:"firmName" json:"firmName" validate:"required"`
	ContactNumber      string     `bson:"contactNumber" json:"contactNumber" validate:"required"`
	GstNumber          string     `bson:"gstNumber" json:"gstNumber" validate:"required,gstNumber"`
	OS                 string     `bson:"os" json:"os" validate:"required"`
	DeviceId           string     `bson:"deviceId" json:"deviceId" binding:"required" validate:"required"`
	DeviceType         DeviceType `bson:"deviceType" json:"deviceType" binding:"required,enum" validate:"required"`
}
