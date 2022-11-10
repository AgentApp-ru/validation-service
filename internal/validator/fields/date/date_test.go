package date

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestGetDateDate(t *testing.T) {
	minField := `"2018-01-01"`
	rawMinValue := json.RawMessage(minField)

	expectedDate, _ := time.Parse("2006-01-02", "2018-01-01")

	if date, ok := getPlainDate(rawMinValue); !ok || !date.Equal(expectedDate) {
		t.Errorf("%v should be equal to %v", date, expectedDate)
	}
}

func CreateFixture(types, birthDate string) (DateDependingConditionFormulaValue, *sync.Map, *sync.Map) {
	var (
		selfMap   *sync.Map
		fieldsMap *sync.Map
		formula   DateDependingConditionFormulaValue
		now       time.Time
	)
	switch types {
	case "min":
		formula = DateDependingConditionFormulaValue{
			Dependency: dependency{
				Value: json.RawMessage(`{
							"scope": "self",
							"key": "birth_date"
						  }`),
				Type: "depending"},
			Operation: "add",
			Unit:      "year",
			Value: value{Depending: dependency{Value: nil,
				Type: "now"},
				Type:      "diff-intervals",
				Direction: "lte",
				Intervals: []intervals{
					{Diff: 20, Value: 14},
					{Diff: 45, Value: 20},
				},
				Unit:    "year",
				Default: 45,
				Offset:  offset{Value: 90, Unit: "day"},
			},
		}
	case "max":
		formula = DateDependingConditionFormulaValue{
			Dependency: dependency{
				Value: json.RawMessage(`{
							"scope": "self",
							"key": "birth_date"
						  }`),
				Type: "depending"},
			Operation: "add",
			Unit:      "year",
			Value: value{Depending: dependency{Value: nil,
				Type: "now"},
				Type:      "diff-intervals",
				Direction: "gte",
				Intervals: []intervals{
					{Diff: 45, Value: 100},
					{Diff: 20, Value: 45},
				},
				Unit:    "year",
				Default: 20,
				Offset:  offset{Value: 90, Unit: "day"},
			},
		}
	}
	selfMap = new(sync.Map)
	fieldsMap = new(sync.Map)
	selfMap.Store("birth_date", birthDate)
	fieldsMap.Store("birth_date", birthDate)

	now = time.Date(2022, 11, 01, 0, 0, 0, 0, time.UTC)
	//mock time.Now
	CurrentDate = func() time.Time {
		return now
	}

	return formula, selfMap, fieldsMap
}

func CalculateDate(t *testing.T, dependingDate string) (time.Time, time.Time) {
	var (
		minDate time.Time
		maxDate time.Time
		ok      bool
	)
	formula, selfMap, fieldsMap := CreateFixture("min", dependingDate)
	if minDate, ok = getRangeExpectedDate(formula, selfMap, fieldsMap); !ok {
		t.Error("не получена дата")
	}

	formula, selfMap, fieldsMap = CreateFixture("max", dependingDate)
	if maxDate, ok = getRangeExpectedDate(formula, selfMap, fieldsMap); !ok {
		t.Error("не получена дата")
	}
	return minDate, maxDate
}

func TestDiffIntervalsMiddle(t *testing.T) {
	dependedDate := "1994-04-24"
	rawDate := "2018-01-01"

	minDate, maxDate := CalculateDate(t, dependedDate)

	expectedDate, _ := time.Parse("2006-01-02", rawDate)
	if minDate.After(expectedDate) || maxDate.Before(expectedDate) {
		t.Errorf("Ошибка %s < %s < %s", minDate, expectedDate, maxDate)
	}
}

func TestDiffIntervalsStart(t *testing.T) {
	dependedDate := "2008-04-24"
	rawDate := "2022-04-24"

	minDate, maxDate := CalculateDate(t, dependedDate)

	expectedDate, _ := time.Parse("2006-01-02", rawDate)
	if minDate.After(expectedDate) || maxDate.Before(expectedDate) {
		t.Errorf("Ошибка %s < %s < %s", minDate, expectedDate, maxDate)
	}
}

func TestDiffIntervalsEnd(t *testing.T) {
	dependedDate := "1950-04-24"
	rawDate := "2022-04-24"

	minDate, maxDate := CalculateDate(t, dependedDate)

	expectedDate, _ := time.Parse("2006-01-02", rawDate)
	if minDate.After(expectedDate) || maxDate.Before(expectedDate) {
		t.Errorf("Ошибка %s < %s < %s", minDate, expectedDate, maxDate)
	}
}

func TestDiffIntervals_dependingDate_EqualDayAndMonth_rawDateErrorStart(t *testing.T) {
	dependedDate := "2008-11-01"
	rawDate := "2022-11-01"

	minDate, maxDate := CalculateDate(t, dependedDate)

	expectedDate, _ := time.Parse("2006-01-02", rawDate)
	if minDate.After(expectedDate) && maxDate.Before(expectedDate) {
		t.Errorf("Ошибка %s < %s < %s", minDate, expectedDate, maxDate)
	}
}

func TestDiffIntervals_dependingDate_EqualDayAndMonth_rawDateOKStart(t *testing.T) {
	dependedDate := "2008-10-31"
	rawDate := "2022-11-01"

	minDate, maxDate := CalculateDate(t, dependedDate)

	expectedDate, _ := time.Parse("2006-01-02", rawDate)
	if minDate.After(expectedDate) && maxDate.Before(expectedDate) {
		t.Errorf("Ошибка %s < %s < %s", minDate, expectedDate, maxDate)
	}
}
