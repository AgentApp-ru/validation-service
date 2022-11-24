package date

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type (
	dependency struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	DateDependingFormulaValue struct {
		Dependency dependency `json:"depending"`
		Operation  string     `json:"operation"`
		Value      int        `json:"value"`
		Unit       string     `json:"unit"`
	}

	intervals struct {
		Diff  int `json:"diff"`
		Value int `json:"value"`
	}

	value struct {
		Depending dependency  `json:"depending"`
		Type      string      `json:"type"`
		Direction string      `json:"direction"`
		Intervals []intervals `json:"intervals"`
		Unit      string      `json:"unit"`
		Default   int         `json:"default"`
	}

	DateDependingConditionFormulaValue struct {
		Dependency dependency `json:"depending"`
		Operation  string     `json:"operation"`
		Unit       string     `json:"unit"`
		Value      value      `json:"value"`
	}
)

func (d *dependency) getInitialDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	switch d.Type {
	case "now":
		return time.Now(), nil
	case "depending":
		var dependingScope DependingValue
		err := json.Unmarshal(d.Value, &dependingScope)
		if err != nil {
			return time.Time{}, err
		}

		value := waitingForValue(dependingScope.Scope, dependingScope.Key, selfMap, fieldsMap)
		if value == nil {
			return time.Time{}, fmt.Errorf("depending field not found: %s/%s", dependingScope.Scope, dependingScope.Key)
		}
		expectedDateRaw := value.(string)
		expectedDate, err := time.Parse("2006-01-02", expectedDateRaw)
		return expectedDate, err
	case "depending_formula":
		var dv = DateDependingFormulaValue{}
		err := json.Unmarshal(d.Value, &dv)
		if err != nil {
			return time.Time{}, err
		}

		expectedDate, err := dv.getExpectedDate(selfMap, fieldsMap)
		if err != nil {
			return time.Time{}, err
		}
		return expectedDate, err
	default:
		return time.Time{}, fmt.Errorf("no logic for dependency: %v", d.Type)
	}
}

func (f *DateDependingFormulaValue) getExpectedDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	switch f.Unit {
	case "year":
		years = f.Value
	case "month":
		months = f.Value
	case "day":
		days = f.Value
	}

	switch f.Operation {
	case "subtract":
		years = 0 - years
		months = 0 - months
		days = 0 - days
	case "add":
		years = 0 + years
		months = 0 + months
		days = 0 + days
	}

	return initialDate.AddDate(years, months, days), nil
}

func waitingForValue(scope, key string, selfMap, fieldsMap *sync.Map) interface{} {
	var scopeObjectMap *sync.Map

	if scope == "self" {
		scopeObjectMap = selfMap
	} else {
		scopeObjectMapLoaded, ok := fieldsMap.Load(scope)
		if !ok {
			return nil
		}
		scopeObjectMap = scopeObjectMapLoaded.(*sync.Map)
	}

	if value, ok := scopeObjectMap.Load(key); ok {
		return value
	}

	return nil
}
