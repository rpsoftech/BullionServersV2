package interfaces

type bullionGeneralUserConfig struct {
	AutoApprove bool `bson:"autoApprove" json:"autoApprove"`
	AutoLogin   bool `bson:"autoLogin" json:"autoLogin"`
}

type BullionSiteInfo struct {
	BaseEntity      `bson:"inline"`
	Name            string                    `bson:"name" json:"name"`
	ShortName       string                    `bson:"shortName" json:"shortName"`
	Domain          string                    `bson:"domain" json:"domain"`
	GeneralUserInfo *bullionGeneralUserConfig `bson:"generalUserInfo" json:"generalUserInfo"`
}

func (b *BullionSiteInfo) AddGeneralUserInfo(AutoApprove bool, AutoLogin bool) (r *BullionSiteInfo) {
	b.GeneralUserInfo = &bullionGeneralUserConfig{
		AutoApprove: AutoApprove,
		AutoLogin:   AutoLogin,
	}
	return b
}

func CreateNewBullionSiteInfo(name string, domain string) (r *BullionSiteInfo) {
	b := BullionSiteInfo{
		BaseEntity: BaseEntity{},
	}
	b.Name = name
	b.Domain = domain
	b.CreateNewId()
	return &b
}
