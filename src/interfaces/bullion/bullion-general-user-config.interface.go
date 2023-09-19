package bullion

type bullionGeneralUserConfig struct {
	AutoApprove bool `bson:"autoApprove"`
	AutoLogin   bool `bson:"autoLogin"`
}
