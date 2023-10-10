package interfaces

import "github.com/rpsoftech/bullion-server/src/validator"

type CalculationSymbolEnum string

const (
	CALCULATION_SYMBOL_GOLD       CalculationSymbolEnum = "GOLD"
	CALCULATION_SYMBOL_SILVER     CalculationSymbolEnum = "SILVER"
	CALCULATION_SYMBOL_GOLD_MCX   CalculationSymbolEnum = "GOLD_MCX"
	CALCULATION_SYMBOL_SILVER_MCX CalculationSymbolEnum = "SILVER_MCX"
)

var (
	calculationSymbolEnumMap = map[string]CalculationSymbolEnum{
		"GOLD":       CALCULATION_SYMBOL_GOLD,
		"SILVER":     CALCULATION_SYMBOL_SILVER,
		"GOLD_MCX":   CALCULATION_SYMBOL_GOLD_MCX,
		"SILVER_MCX": CALCULATION_SYMBOL_SILVER_MCX,
	}
)

func init() {
	validator.RegisterEnumValidatorFunc("CalculationSymbolEnum", ValidateEnumCalculationSymbolEnum)
}

func ValidateEnumCalculationSymbolEnum(value string) bool {
	_, ok := calculationSymbolEnumMap[value]
	return ok
}

func (s CalculationSymbolEnum) String() string {
	switch s {
	case CALCULATION_SYMBOL_GOLD:
		return "GOLD"
	case CALCULATION_SYMBOL_SILVER:
		return "SILVER"
	case CALCULATION_SYMBOL_GOLD_MCX:
		return "GOLD_MCX"
	case CALCULATION_SYMBOL_SILVER_MCX:
		return "SILVER_MCX"
	}
	return "unknown"
}

func (s CalculationSymbolEnum) IsValid() bool {
	switch s {
	case
		CALCULATION_SYMBOL_GOLD,
		CALCULATION_SYMBOL_SILVER,
		CALCULATION_SYMBOL_GOLD_MCX,
		CALCULATION_SYMBOL_SILVER_MCX:
		return true
	}

	return false
}
