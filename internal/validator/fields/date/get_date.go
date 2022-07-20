package date

import (
	"encoding/json"
	"sync"
	"time"
	"validation_service/pkg/log"
)

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

func getRangeExpectedDate(formula DateDependingConditionFormulaValue, selfMap, fieldsMap *sync.Map) (time.Time, bool) {
	var (
		today  time.Time
		result expectedDate
	)

	today = time.Now()
	result.Value = today

	for _, item := range formula.Condition.Items {
		dv := DateDependingFormulaValue{
			formula.Dependency,
			formula.Operation,
			item,
			formula.Unit}
		date, err := dv.getExpectedDate(selfMap, fieldsMap)
		if err != nil {
			log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
			return time.Time{}, false
		}

		switch formula.Direction {
		case "gte":
			if date.Before(today) {
				var item = expectedDate{Value: date, Sub: int(today.Sub(date))}
				if result.Sub == 0 || item.Sub < result.Sub {
					result = item
				}
			}
		case "lte":
			if date.After(today) {
				var item = expectedDate{Value: date, Sub: int(today.Sub(date))}
				if item.Sub > result.Sub {
					result = item
				}
			}
		}
	}
	return result.Value, true
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

	switch formula.Condition.Type {
	case "range":
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
		formula.Condition.Default,
		formula.Unit}

	date, err := dv.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return time.Time{}, false
	}
	return date, true
}
