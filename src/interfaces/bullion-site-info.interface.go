package interfaces

type bullionGeneralUserConfig struct {
	AutoApprove bool `bson:"autoApprove" json:"autoApprove" validate:"boolean"`
	AutoLogin   bool `bson:"autoLogin" json:"autoLogin" validate:"boolean"`
}
type bullionConfigs struct {
	OTPLength               int  `bson:"otpLength" json:"otpLength" validate:"required"`
	HaveCustomWhatsappAgent bool `bson:"haveCustomWhatsappAgent" json:"haveCustomWhatsappAgent" validate:"boolean"`
	// AutoApprove bool `bson:"autoApprove" json:"autoApprove" validate:"boolean"`
	// AutoLogin   bool `bson:"autoLogin" json:"autoLogin" validate:"boolean"`
}

type BullionSiteBasicInfo struct {
	Name      string `bson:"name" json:"name" validate:"required"`
	ShortName string `bson:"shortName" json:"shortName" validate:"required"`
	Domain    string `bson:"domain" json:"domain" validate:"required"`
}
type BullionSiteInfoEntity struct {
	*BaseEntity           `bson:"inline"`
	*BullionSiteBasicInfo `bson:"inline"`
	BullionConfigs        *bullionConfigs           `bson:"bullionConfigs" json:"-" validate:"required"`
	GeneralUserInfo       *bullionGeneralUserConfig `bson:"generalUserInfo" json:"-" validate:"required"`
}

func (b *BullionSiteInfoEntity) AddGeneralUserInfo(AutoApprove bool, AutoLogin bool) (r *BullionSiteInfoEntity) {
	b.GeneralUserInfo = &bullionGeneralUserConfig{
		AutoApprove: AutoApprove,
		AutoLogin:   AutoLogin,
	}
	return b
}

func CreateNewBullionSiteInfo(name string, shortName string, domain string) *BullionSiteInfoEntity {
	b := BullionSiteInfoEntity{
		BaseEntity: &BaseEntity{},
		BullionSiteBasicInfo: &BullionSiteBasicInfo{
			Name:      name,
			ShortName: shortName,
			Domain:    domain,
		},
	}
	b.createNewId()
	return &b
}
