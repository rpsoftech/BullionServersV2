package interfaces

import "github.com/rpsoftech/bullion-server/src/validator"

type CalculateOnPriceType string

const (
	CALCULATE_ON_BID_ASK CalculateOnPriceType = "BID_ASK"
	CALCULATE_ON_BID     CalculateOnPriceType = "BID"
	CALCULATE_ON_ASK     CalculateOnPriceType = "ASK"
)

var (
	calculateOnPriceTypeMap = EnumValidatorBase{
		Data: map[string]interface{}{
			"BID_ASK": CALCULATE_ON_BID_ASK,
			"ASK":     CALCULATE_ON_ASK,
			"BID":     CALCULATE_ON_BID,
		},
	}
)

func init() {
	validator.RegisterEnumValidatorFunc("CalculateOnPriceType", calculateOnPriceTypeMap.Validate)
}

// ValidateEnumCalculateOnPriceType checks if the given value is a valid calculateOnPriceType.
//
// value: the value to be validated as a calculateOnPriceType.
// Returns: true if the value is a valid calculateOnPriceType, false otherwise.

func (s CalculateOnPriceType) String() string {
	switch s {
	case CALCULATE_ON_ASK:
		return "ASK"
	case CALCULATE_ON_BID:
		return "BID"
	case CALCULATE_ON_BID_ASK:
		return "BID_ASK"
	}
	return "unknown"
}

func (s CalculateOnPriceType) IsValid() bool {
	switch s {
	case
		CALCULATE_ON_ASK,
		CALCULATE_ON_BID,
		CALCULATE_ON_BID_ASK:
		return true
	}

	return false
}
