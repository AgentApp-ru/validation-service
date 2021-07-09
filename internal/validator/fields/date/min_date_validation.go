package date

import (
	"encoding/json"
	"sync"
	"time"
	"validation_service/pkg/log"
)

func (dv *DateValidator) validateMinDate(fieldDate time.Time, pattern *DatePattern) bool {
	var (
		expectedDate time.Time
		ok           bool
	)

	for _, minPattern := range pattern.Min {
		switch minPattern.PatternType {
		case "date":
			expectedDate, ok = getMinDate(fieldDate, minPattern.Value)
		case "now":
			expectedDate, ok = getMinDateNow(fieldDate, minPattern.Value)
		case "depending":
			expectedDate, ok = getMinDepending(fieldDate, minPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_formula":
			expectedDate, ok = getMinDependingFormula(fieldDate, minPattern.Value, dv.objectMap, dv.allFieldsMap)
		case "depending_condition_formula":
			expectedDate, ok = getMinDependingConditionFormula(fieldDate, minPattern.Value, dv.objectMap, dv.allFieldsMap)
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
			ok = false
		}
	}
	if !ok {
		return false
	}

	return !fieldDate.Before(expectedDate)
}

func getMinDate(fieldDate time.Time, rawValue json.RawMessage) (time.Time, bool) {
	var value DateDateValue
	err := json.Unmarshal(rawValue, &value)
	if err != nil {
		// TODO: log error
		return time.Time{}, false
	}

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		// TODO: log error
		return time.Time{}, false
	}

	return expectedDate, true
}

func getMinDateNow(fieldDate time.Time, rawValue json.RawMessage) (time.Time, bool) {
	return time.Now(), true
}

func getMinDepending(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var value DependingValue
	json.Unmarshal(rawValue, &value)

	dependingValue := waitingForValue(value.Scope, value.Key, selfMap, fieldsMap)
	if dependingValue == nil {
		log.Logger.Errorf("depending field not found: %s/%s", value.Scope, value.Key)
		return time.Time{}, false
	}
	expectedDateRaw := dependingValue.(string)
	expectedDate, err := time.Parse("2006-01-02", expectedDateRaw)
	if err != nil {
		return time.Time{}, false
	}

	return expectedDate, true
}

func getMinDependingFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}

	return expectedDate, true
}

func getMinDependingConditionFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var formula DateDependingConditionFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}

	return expectedDate, true
}
