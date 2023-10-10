package interfaces

import "github.com/rpsoftech/bullion-server/src/validator"

type BaseSymbolEnum string

const (
	BASE_SYMBOL_GOLD   BaseSymbolEnum = "GOLD"
	BASE_SYMBOL_SILVER BaseSymbolEnum = "SILVER"
)

var (
	baseSymbolEnumMap = map[string]BaseSymbolEnum{
		"GOLD":   BASE_SYMBOL_GOLD,
		"SILVER": BASE_SYMBOL_SILVER,
	}
)

func init() {
	validator.RegisterEnumValidatorFunc("BaseSymbolEnum", ValidateBaseSymbolEnum)
}

func ValidateBaseSymbolEnum(value string) bool {
	_, ok := baseSymbolEnumMap[value]
	return ok
}

func (s BaseSymbolEnum) String() string {
	switch s {
	case BASE_SYMBOL_GOLD:
		return "GOLD"
	case BASE_SYMBOL_SILVER:
		return "SILVER"
	}
	return "unknown"
}

func (s BaseSymbolEnum) IsValid() bool {
	switch s {
	case
		BASE_SYMBOL_GOLD,
		BASE_SYMBOL_SILVER:
		return true
	}

	return false
}
