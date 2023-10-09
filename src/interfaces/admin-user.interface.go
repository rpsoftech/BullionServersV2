package interfaces

type AdminUser struct {
	*BaseEntity         `bson:"inline"`
	*UserRolesInterface `bson:"inline"`
	UserName            string `bson:"userName" json:"userName" validate:"require"`
	Password            string `bson:"password" json:"password" validate:"require"`
	NickName            string `bson:"nickName" json:"nickName" validate:"require"`
	BullionId           string `bson:"bullionId" json:"bullionId" validate:"require"`
}

func (admin *AdminUser) MatchPassword(password string) bool {
	return admin.Password == password
}

func (admin *AdminUser) CreateNewEntity(UserName string, Password string, NickName string, BullionId string) *AdminUser {
	admin.BaseEntity = &BaseEntity{}
	admin.UserName = UserName
	admin.Password = Password
	admin.NickName = NickName
	admin.BullionId = BullionId
	admin.UserRolesInterface = &UserRolesInterface{
		Role: ROLE_ADMIN,
	}
	admin.createNewId()
	return admin
}
