package interfaces

import "github.com/rpsoftech/bullion-server/src/validator"

type SymbolsEnum string

const (
	SYMBOL_GOLD        SymbolsEnum = "GOLD"
	SYMBOL_SILVER      SymbolsEnum = "SILVER"
	SYMBOL_GOLD_MCX    SymbolsEnum = "GOLD_MCX"
	SYMBOL_SILVER_MCX  SymbolsEnum = "SILVER_MCX"
	SYMBOL_GOLD_NEXT   SymbolsEnum = "GOLD_NEXT"
	SYMBOL_SILVER_NEXT SymbolsEnum = "SILVER_NEXT"
	SYMBOL_GOLD_SPOT   SymbolsEnum = "GOLD_SPOT"
	SYMBOL_SILVER_SPOT SymbolsEnum = "SILVER_SPOT"
	SYMBOL_INR         SymbolsEnum = "INR"
)

var (
	symbolEnumMap = EnumValidatorBase{
		Data: map[string]interface{}{
			"GOLD":        SYMBOL_GOLD,
			"GOLD_MCX":    SYMBOL_GOLD_MCX,
			"GOLD_SPOT":   SYMBOL_GOLD_SPOT,
			"GOLD_NEXT":   SYMBOL_GOLD_NEXT,
			"SILVER":      SYMBOL_SILVER,
			"SILVER_MCX":  SYMBOL_SILVER_MCX,
			"SILVER_NEXT": SYMBOL_SILVER_NEXT,
			"SILVER_SPOT": SYMBOL_SILVER_SPOT,
			"INR":         SYMBOL_INR,
		},
	}
)

func init() {
	validator.RegisterEnumValidatorFunc("SymbolsEnum", symbolEnumMap.Validate)
}
func (s SymbolsEnum) String() string {
	switch s {
	case SYMBOL_GOLD:
		return "GOLD"
	case SYMBOL_SILVER:
		return "SILVER"
	case SYMBOL_GOLD_MCX:
		return "GOLD_MCX"
	case SYMBOL_SILVER_MCX:
		return "SILVER_MCX"
	case SYMBOL_GOLD_NEXT:
		return "GOLD_NEXT"
	case SYMBOL_SILVER_NEXT:
		return "SILVER_NEXT"
	case SYMBOL_GOLD_SPOT:
		return "GOLD_SPOT"
	case SYMBOL_SILVER_SPOT:
		return "SILVER_SPOT"
	case SYMBOL_INR:
		return "INR"

	}
	return "unknown"
}

func (s SymbolsEnum) IsValid() bool {
	switch s {
	case
		SYMBOL_GOLD,
		SYMBOL_SILVER,
		SYMBOL_GOLD_MCX,
		SYMBOL_SILVER_MCX,
		SYMBOL_GOLD_NEXT,
		SYMBOL_SILVER_NEXT,
		SYMBOL_GOLD_SPOT,
		SYMBOL_SILVER_SPOT,
		SYMBOL_INR:
		return true
	}

	return false
}
