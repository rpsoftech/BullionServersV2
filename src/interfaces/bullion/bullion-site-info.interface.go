package bullion

import "github.com/rpsoftech/bullion-server/src/interfaces"

type BullionSiteInfo struct {
	*interfaces.BaseEntity `bson:"inline"`
	Name                   string                    `bson:"name"`
	Domains                []string                  `bson:"domains"`
	GeneralUserInfo        *bullionGeneralUserConfig `bson:"generalUserInfo"`
}

func (b *BullionSiteInfo) AddGeneralUserInfo(AutoApprove bool, AutoLogin bool) (r *BullionSiteInfo) {
	b.GeneralUserInfo = &bullionGeneralUserConfig{
		AutoApprove: AutoApprove,
		AutoLogin:   AutoLogin,
	}
	return b
}
