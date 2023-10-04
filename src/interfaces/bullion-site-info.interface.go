package interfaces

type bullionGeneralUserConfig struct {
	AutoApprove bool `bson:"autoApprove" json:"autoApprove" validate:"required"`
	AutoLogin   bool `bson:"autoLogin" json:"autoLogin" validate:"required"`
}

type BullionSiteInfoEntity struct {
	BaseEntity      `bson:"inline"`
	Name            string                    `bson:"name" json:"name" validate:"required"`
	ShortName       string                    `bson:"shortName" json:"shortName" validate:"required"`
	Domain          string                    `bson:"domain" json:"domain" validate:"required"`
	GeneralUserInfo *bullionGeneralUserConfig `bson:"generalUserInfo" json:"generalUserInfo" validate:"required"`
}

func (b *BullionSiteInfoEntity) AddGeneralUserInfo(AutoApprove bool, AutoLogin bool) (r *BullionSiteInfoEntity) {
	b.GeneralUserInfo = &bullionGeneralUserConfig{
		AutoApprove: AutoApprove,
		AutoLogin:   AutoLogin,
	}
	return b
}

func CreateNewBullionSiteInfo(name string, domain string) (r *BullionSiteInfoEntity) {
	b := BullionSiteInfoEntity{
		BaseEntity: BaseEntity{},
	}
	b.Name = name
	b.Domain = domain
	b.CreateNewId()
	return &b
}
