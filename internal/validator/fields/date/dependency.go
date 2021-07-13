package date

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"validation_service/pkg/log"
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

	condition struct {
		Type    string         `json:"type"`
		Items   map[string]int `json:"items"`
		Default int            `json:"default"`
	}

	conditionValue struct {
		Field     DependingValue `json:"field"`
		Condition condition      `json:"condition"`
	}

	DateDependingConditionFormulaValue struct {
		Dependency     dependency     `json:"depending"`
		Operation      string         `json:"operation"`
		ConditionValue conditionValue `json:"value"`
		Unit           string         `json:"unit"`
	}
)

func (c *condition) getItem(searchedValue string) (int, error) {
	switch c.Type {
	case "equals":
		item, ok := c.Items[searchedValue]
		if !ok {
			return c.Default, nil
		}
		return item, nil
	default:
		log.Logger.Errorf("unknown condition type: %s", c.Type)
		return c.Default, nil
	}
}

func (c *conditionValue) getItem(selfMap, fieldsMap *sync.Map) (int, error) {
	value := waitingForValue(c.Field.Scope, c.Field.Key, selfMap, fieldsMap)
	if value == nil {
		value = "default"
	}
	searchedValue := fmt.Sprintf("%v", value)

	return c.Condition.getItem(searchedValue)
}

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

	if days == 0 {
		days = -1
	}

	return initialDate.AddDate(years, months, days), nil
}

func (f *DateDependingConditionFormulaValue) getExpectedDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	item, err := f.ConditionValue.getItem(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	switch f.Unit {
	case "year":
		years = item
	case "month":
		months = item
	case "day":
		days = item
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

	if days == 0 {
		days = -1
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

	end := time.Now().Add(10 * time.Second)
	for time.Now().Before(end) {
		if value, ok := scopeObjectMap.Load(key); ok {
			return value
		}
		time.Sleep(300 * time.Millisecond)
	}

	return nil
}
