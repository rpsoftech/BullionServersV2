package bullion

import "github.com/rpsoftech/bullion-server/src/interfaces"

type BullionSiteInfo struct {
	*interfaces.BaseEntity `bson:"inline"`
	Name                   string                    `bson:"name"`
	Domain                 string                    `bson:"domain"`
	GeneralUserInfo        *bullionGeneralUserConfig `bson:"generalUserInfo"`
}

func (b *BullionSiteInfo) AddGeneralUserInfo(AutoApprove bool, AutoLogin bool) (r *BullionSiteInfo) {
	b.GeneralUserInfo = &bullionGeneralUserConfig{
		AutoApprove: AutoApprove,
		AutoLogin:   AutoLogin,
	}
	return b
}

func CreateNewBullionSiteInfo(name string, domain string) (r *BullionSiteInfo) {
	b := &BullionSiteInfo{
		BaseEntity: &interfaces.BaseEntity{},
	}
	b.Name = name
	b.Domain = domain
	b.CreateNewId()
	return b
}
