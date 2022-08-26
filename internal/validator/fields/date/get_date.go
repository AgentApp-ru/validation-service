package date

import (
	"encoding/json"
	"sync"
	"time"
	"validation_service/pkg/log"
)

const secondsInYear = 31536000

type (
	expectedDate struct {
		Value time.Time
		Sub   int
	}
)

func getPlainDate(rawValue json.RawMessage) (time.Time, bool) {
	var value DateDateValue
	err := json.Unmarshal(rawValue, &value)
	if err != nil {
		// TODO: log error
		println("!!!", err.Error())
		return time.Time{}, false
	}

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		// TODO: log error
		println(err)
		return time.Time{}, false
	}

	return expectedDate, true
}

func getPlainNow(rawValue json.RawMessage) (time.Time, bool) {
	return time.Now().Truncate(24 * time.Hour), true
}

func getDepending(rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var value DependingValue
	err := json.Unmarshal(rawValue, &value)
	if err != nil {
		log.Logger.Errorf("Ошибка парсинга JSON: %s", rawValue)
		return time.Time{}, false
	}

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

func getDependingFormula(rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var formula DateDependingFormulaValue
	err := json.Unmarshal(rawValue, &formula)
	if err != nil {
		log.Logger.Errorf("Ошибка парсинга JSON: %s", rawValue)
		return time.Time{}, false
	}

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}

	return expectedDate, true
}

func _getExpectedDate(selfMap, fieldsMap *sync.Map, formula DateDependingConditionFormulaValue, value int) time.Time {
	dv := DateDependingFormulaValue{
		formula.Dependency,
		formula.Operation,
		value,
		formula.Unit}
	date, err := dv.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}
	}
	return date
}

func getRangeExpectedDate(formula DateDependingConditionFormulaValue, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var result expectedDate
	var dependingDate = _getDependingDate(formula)
	var defaultDate = _getExpectedDate(selfMap, fieldsMap, formula, formula.Value.Default)

	for _, inter := range formula.Value.Intervals {
		date := _getExpectedDate(selfMap, fieldsMap, formula, inter.Value)
		result = expectedDate{Value: dependingDate, Sub: int(dependingDate.Sub(date).Seconds() / secondsInYear)}
		switch formula.Value.Direction {
		case "gte":
			if date.Before(dependingDate) {
				continue
			}
			if result.Sub >= inter.Diff {
				date = defaultDate
			}
			result = expectedDate{Value: date, Sub: int(date.Sub(dependingDate).Seconds() / secondsInYear)}

		case "lte":
			if date.After(dependingDate) {
				continue
			}
			if result.Sub >= inter.Diff {
				date = defaultDate
			}
			result = expectedDate{Value: date, Sub: int(dependingDate.Sub(date).Seconds() / secondsInYear)}
		}
	}
	if result.Value.After(dependingDate) {
		result.Value = dependingDate
	}

	return result.Value, true
}

func _getDependingDate(formula DateDependingConditionFormulaValue) time.Time {
	switch formula.Value.Depending.Type {
	case "now":
		return time.Now()
	default:
	}

	return time.Now()
}

func getDependingConditionFormula(rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var (
		formula      DateDependingConditionFormulaValue
		ExpectedDate time.Time
		ok           bool
	)

	err := json.Unmarshal(rawValue, &formula)
	if err != nil {
		log.Logger.Errorf("Ошибка парсинга JSON: %s", rawValue)
		return time.Time{}, false
	}

	switch formula.Value.Type {
	case "diff-intervals":
		ExpectedDate, ok = getRangeExpectedDate(formula, selfMap, fieldsMap)
	default:
		ExpectedDate, ok = getDefaultExpectedDate(formula, selfMap, fieldsMap)
	}

	return ExpectedDate, ok
}

func getDefaultExpectedDate(formula DateDependingConditionFormulaValue, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	dv := DateDependingFormulaValue{
		formula.Dependency,
		formula.Operation,
		formula.Value.Default,
		formula.Unit}

	date, err := dv.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}
	return date, true
}
