package date

import (
	"encoding/json"
	"sync"
	"time"
	"validation_service/pkg/log"
)

const secondsInYear = 31536000
const secondsInDay = 86400

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
	var (
		dependingValueDateWithoutOffset time.Time
		ok                              bool
		passed                          bool
		diff                            int
	)
	result := _getExpectedDate(selfMap, fieldsMap, formula, formula.Value.Default)
	dependingDate := _getDependingDate(formula, selfMap, fieldsMap)
	dependingValueDate, _ := formula.Dependency.getInitialDate(selfMap, fieldsMap)

	for _, inter := range formula.Value.Intervals {
		switch formula.Value.Direction {
		case "gte":
			direction := -1
			if dependingValueDateWithoutOffset, diff, ok, _ = checkOffsetDateValue(formula, dependingValueDate, dependingDate, direction); ok {
				return dependingValueDateWithoutOffset, ok
			}
			if diff >= inter.Diff {
				newResult := _getExpectedDate(selfMap, fieldsMap, formula, inter.Value)
				if result.Before(newResult) {
					result = newResult
				}
			}

		case "lte":
			direction := 1
			if dependingValueDateWithoutOffset, diff, ok, passed = checkOffsetDateValue(formula, dependingValueDate, dependingDate, direction); ok {
				return dependingValueDateWithoutOffset, ok
			}
			if passed {
				break
			}
			if diff <= inter.Diff {
				newResult := _getExpectedDate(selfMap, fieldsMap, formula, inter.Value)
				if result.After(newResult) {
					result = newResult
				}
			}
		}
	}

	if result.After(dependingDate) {
		result = CurrentDate()
	}

	return result, true
}

func checkOffsetDateValue(formula DateDependingConditionFormulaValue, dependingValueDate, dependingDate time.Time, direction int) (time.Time, int, bool, bool) {
	//check birthday
	var (
		dv                              DateDependingFormulaValue
		dependingValueDateWithoutOffset time.Time
		diff                            int
	)
	err := json.Unmarshal(formula.Dependency.Value, &dv)
	if err != nil {
		return time.Time{}, 0, false, false
	}

	dependingValueDateWithoutOffset = dependingValueDate.AddDate(0, 0, dv.Value*direction)
	diff = int(dependingDate.Sub(dependingValueDateWithoutOffset).Seconds()) / secondsInYear
	for _, inter := range formula.Value.Intervals {
		if inter.Diff == diff || inter.Value == diff {
			if dependingValueDateWithoutOffset.Day() == dependingDate.Day() && dependingValueDateWithoutOffset.Month() == dependingDate.Month() {
				return dependingDate.AddDate(0, 0, direction), diff, true, false
			}

			offsetDate := time.Date(
				dependingDate.Year(),
				dependingValueDateWithoutOffset.Month(),
				dependingValueDateWithoutOffset.Day(),
				0, 0, 0, 0, time.UTC,
			)
			offsetDays := int(dependingDate.Sub(offsetDate).Seconds()) / secondsInDay

			if offsetDays > dv.Value {
				return time.Time{}, diff, false, true
			}
		}
	}

	return time.Time{}, diff, false, false
}

func _getDependingDate(formula DateDependingConditionFormulaValue, selfMap, fieldsMap *sync.Map) time.Time {
	switch formula.Value.Depending.Type {
	case "now":
		return CurrentDate()
	case "depending_formula":
		var dv = DateDependingFormulaValue{}
		err := json.Unmarshal(formula.Value.Depending.Value, &dv)
		if err != nil {
			log.Logger.Errorf("Ошибка парсинга JSON: %s", formula.Dependency.Value)
			return time.Now()
		}

		expectedDate, err := dv.getExpectedDate(selfMap, fieldsMap)
		if err != nil {
			log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
			return time.Now()
		}
		return expectedDate
	default:
	}

	return time.Now()
}

var CurrentDate = func() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
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
