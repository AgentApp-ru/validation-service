package date

import (
	"time"
	"validation_service/pkg/log"
)

func (dv *DateValidator) validateMinDate(fieldDate time.Time, pattern *DatePattern) bool {
	var (
		expectedDate time.Time
		ok           bool
	)

	if len(pattern.Min) == 0 {
		return true
	}

	for _, minPattern := range pattern.Min {
		switch minPattern.PatternType {
		case "date":
			expectedDate, ok = getPlainDate(minPattern.Value)
		case "now":
			expectedDate, ok = getPlainNow(minPattern.Value)
		case "depending":
			expectedDate, ok = getDepending(minPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_formula":
			expectedDate, ok = getDependingFormula(minPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_condition_formula":
			expectedDate, ok = getDependingConditionFormula(minPattern.Value, dv.objectMap, dv.allFieldsMap)
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
			ok = false
		}
		if !ok || fieldDate.Before(expectedDate) {
			return false
		}
	}

	return true
}
