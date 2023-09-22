package interfaces

type UserRoles string

const (
	ROLE_RATE_ADMIN   UserRoles = "RATE_ADMIN"
	ROLE_SUPER_ADMIN  UserRoles = "SUPER_ADMIN"
	ROLE_ADMIN        UserRoles = "ADMIN"
	ROLE_GENERAL_USER UserRoles = "GENERAL_USER"
	ROLE_TRADE_USER   UserRoles = "TRADE_USER"
	ROLE_GOD          UserRoles = "GOD"
)

func (s UserRoles) String() string {
	switch s {
	case ROLE_RATE_ADMIN:
		return "RATE_ADMIN"
	case ROLE_SUPER_ADMIN:
		return "SUPER_ADMIN"
	case ROLE_ADMIN:
		return "ADMIN"
	case ROLE_GENERAL_USER:
		return "GENERAL_USER"
	case ROLE_TRADE_USER:
		return "TRADE_USER"
	case ROLE_GOD:
		return "GOD"

	}
	return "unknown"
}

func (s UserRoles) IsValid() bool {
	switch s {
	case
		ROLE_RATE_ADMIN,
		ROLE_SUPER_ADMIN,
		ROLE_ADMIN,
		ROLE_GENERAL_USER,
		ROLE_TRADE_USER,
		ROLE_GOD:
		return true
	}

	return false
}

type UserRolesInterface struct {
	Role UserRoles `bson:"role" json:"role" binding:"required,enum"`
}
