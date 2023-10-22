package interfaces

type (
	BankDetailsBase struct {
		BullionId     string `bson:"bullionId" json:"bullionId"`
		AccountName   string `bson:"accountName" json:"accountName"`
		BankName      string `bson:"bankName" json:"bankName"`
		AccountNumber string `bson:"accountNumber" json:"accountNumber"`
		IFSC          string `bson:"ifsc" json:"ifsc"`
		Sequence      string `bson:"sequence" json:"sequence"`
		BranchName    string `bson:"branchName" json:"branchName"`
	}

	BankDetailsEntity struct {
		*BaseEntity      `bson:"inline"`
		*BankDetailsBase `bson:"inline"`
	}
)

func CreateNewBankDetails(base *BankDetailsBase) *BankDetailsEntity {
	entity := &BankDetailsEntity{
		BaseEntity:      &BaseEntity{},
		BankDetailsBase: base,
	}
	entity.createNewId()
	return entity
}
