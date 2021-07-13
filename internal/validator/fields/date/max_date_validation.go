package date

import (
	"encoding/json"
	"sync"
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
			expectedDate, ok = getMaxDate(fieldDate, maxPattern.Value)
		case "now":
			expectedDate, ok = getMaxNow(fieldDate, maxPattern.Value)
		case "depending_formula":
			expectedDate, ok = getMaxDependingFormula(fieldDate, maxPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_condition_formula":
			expectedDate, ok = getMaxDependingConditionFormula(fieldDate, maxPattern.Value, dv.objectMap, dv.allFieldsMap)
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			ok = false
		}
	}
	if !ok {
		return false
	}

	return !fieldDate.After(expectedDate)
}

func getMaxDate(fieldDate time.Time, rawValue json.RawMessage) (time.Time, bool) {
	var value DateDateValue
	json.Unmarshal(rawValue, &value)

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		return time.Time{}, false
	}

	return expectedDate, true
}

func getMaxNow(fieldDate time.Time, rawValue json.RawMessage) (time.Time, bool) {
	return time.Now(), true
}

func getMaxDependingFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}

	return expectedDate, true
}

func getMaxDependingConditionFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var formula DateDependingConditionFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}

	return expectedDate, true
}
