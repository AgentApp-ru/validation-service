package date

import (
	"time"
	"validation_service/pkg/log"
)

func (dv *DateValidator) validateMaxDate(fieldDate time.Time, pattern *DatePattern) bool {
	var (
		expectedDate time.Time
		ok           bool
	)

	if len(pattern.Max) == 0 {
		return true
	}

	for _, maxPattern := range pattern.Max {
		switch maxPattern.PatternType {
		case "date":
			expectedDate, ok = getPlainDate(maxPattern.Value)
		case "now":
			expectedDate, ok = getPlainNow(maxPattern.Value)
		case "depending_formula":
			expectedDate, ok = getDependingFormula(maxPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_condition_formula":
			expectedDate, ok = getDependingConditionFormula(maxPattern.Value, dv.objectMap, dv.allFieldsMap)
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			ok = false
		}
		if !ok || fieldDate.After(expectedDate) {
			return false
		}
	}

	return true
}
