package requirements

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type (
	isRequired           bool
	requiredAlternatives []string
	dependsOn            string

	dependingValue struct {
		Scope string `json:"scope"`
		Key   string `json:"key"`
	}

	dependingValueWithType struct {
		Scope string `json:"scope"`
		Key   string `json:"key"`
		Type  string `json:"type"`
	}

	dependency struct {
		Type      string          `json:"type"`
		Value     json.RawMessage `json:"value"`
		ValueType string          `json:"value_type"`
	}

	dependingFormulaValue struct {
		Dependency dependency `json:"depending"`
		Operation  string     `json:"operation"`
		Value      dependency `json:"value"`
		Unit       string     `json:"unit"`
	}

	condition struct {
		Type    string          `json:"type"`
		Items   map[string]bool `json:"items"`
		Default bool            `json:"default"`
	}

	conditionRequirement struct {
		Type      string          `json:"type"`
		Value     json.RawMessage `json:"value"`
		Condition condition       `json:"condition"`
	}

	required struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	Requirements struct {
		Required  required  `json:"required"`
		DependsOn dependsOn `json:"depends_on"`
	}
)

func (r *Requirements) CheckRequired(isFieldPresent bool, selfMap *sync.Map) bool {
	switch r.Required.Type {
	case "this":
		var value isRequired
		err := json.Unmarshal(r.Required.Value, &value)
		if err != nil {
			// TODO: log error
			println("!!!", err.Error())
			return false
		}

		return !bool(value) || isFieldPresent
	case "any":
		var value requiredAlternatives
		err := json.Unmarshal(r.Required.Value, &value)
		if err != nil {
			// TODO: log error
			println("!!!", err.Error())
			return false
		}

		return isWaitingForFieldSucceed(value, selfMap, 2*time.Second)
	case "depends_on":
		var value dependsOn
		err := json.Unmarshal(r.Required.Value, &value)
		if err != nil {
			// TODO: log error
			// println("!!!", err.Error())
			return false
		}

		if isWaitingForFieldSucceed([]string{string(value)}, selfMap, 150*time.Millisecond) {
			return isFieldPresent
		} else {
			return true
		}
	case "depending_condition_formula":
		var value conditionRequirement
		err := json.Unmarshal(r.Required.Value, &value)
		if err != nil {
			// TODO: log error
			// println("!!!", err.Error())
			return false
		}
		if checkDependingConditionRequirements(value, selfMap) {
			return isFieldPresent
		} else {
			return true
		}
	default:
		return false
	}
}

func isWaitingForFieldSucceed(fieldsToWait []string, scopeObjectMap *sync.Map, timeout time.Duration) bool {
	end := time.Now().Add(timeout)

	for time.Now().Before(end) {
		for _, key := range fieldsToWait {
			if _, ok := scopeObjectMap.Load(key); ok {
				return true
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false
}

func checkDependingConditionRequirements(conditionRequirement conditionRequirement, scopeObjectMap *sync.Map) bool {
	initialValue, err := getInitalValue(conditionRequirement.Type, conditionRequirement.Value, scopeObjectMap)
	if err != nil {
		println(err.Error())
		return false
	}

	switch conditionRequirement.Condition.Type {
	case "starts_with":
		value := initialValue.(string)
		for k, v := range conditionRequirement.Condition.Items {
			if k == value[:len(k)] {
				return v
			}
		}
		return conditionRequirement.Condition.Default
	case "equals":
		valueInt := initialValue.(int)
		for k, v := range conditionRequirement.Condition.Items {
			keyInt, _ := strconv.Atoi(k)
			if valueInt == keyInt {
				return v
			}
		}
		return conditionRequirement.Condition.Default
	}

	return false
}

func getInitalValue(cType string, cValue json.RawMessage, scopeObjectMap *sync.Map) (interface{}, error) {
	switch cType {
	case "depending_formula":
		return getValueFromDependingFormula(cValue, scopeObjectMap)
	case "depending":
		var dependingValue *dependingValue

		err := json.Unmarshal(cValue, &dependingValue)
		if err != nil {
			return nil, err
		}
		initial_value := waitingForValue(dependingValue.Key, scopeObjectMap)
		if initial_value == nil {
			return nil, fmt.Errorf("empty depending value: %s", dependingValue.Key)
		}
		return initial_value, nil
	}

	return nil, fmt.Errorf("no logic for %s", cType)
}

func getValueFromDependingFormula(value json.RawMessage, scopeObjectMap *sync.Map) (interface{}, error) {
	var dependingFormulaValue *dependingFormulaValue

	err := json.Unmarshal(value, &dependingFormulaValue)
	if err != nil {
		return nil, err
	}

	initialValue, err := getInitalSecondValueForFormula(dependingFormulaValue.Dependency, scopeObjectMap)
	if err != nil {
		return nil, err
	}
	secondValue, err := getInitalSecondValueForFormula(dependingFormulaValue.Value, scopeObjectMap)
	if err != nil {
		return nil, err
	}
	result, err := calculate(
		initialValue,
		secondValue,
		dependingFormulaValue.Dependency.ValueType,
		dependingFormulaValue.Value.ValueType,
		dependingFormulaValue.Operation,
		dependingFormulaValue.Unit,
	)

	return result, err
}

func getInitalSecondValueForFormula(dependency dependency, scopeObjectMap *sync.Map) (interface{}, error) {
	var (
		value interface{}
		err   error
	)
	switch dependency.Type {
	case "now":
		value = time.Now().Truncate(24 * time.Hour).Format("2006-02-01")
	case "plain":
		err := json.Unmarshal(dependency.Value, &value)
		if err != nil {
			return nil, err
		}
	case "depending":
		var dependingValueWithType *dependingValueWithType
		err := json.Unmarshal(dependency.Value, &dependingValueWithType)
		if err != nil {
			return nil, err
		}
		value = waitingForValue(dependingValueWithType.Key, scopeObjectMap)
		if value == nil {
			return nil, fmt.Errorf("!! empty depending value: %s", dependingValueWithType.Key)
		}
	case "depending_formula":
		value, err = getValueFromDependingFormula(dependency.Value, scopeObjectMap)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("no logic for %s", dependency.Type)
	}

	return value, nil
}

func calculate(
	initialValue, secondValue interface{},
	initialValueType, secondValueType string,
	operation, unit string,
) (interface{}, error) {
	var (
		initialValueDate time.Time
		secondValueDate  time.Time
		secondValueInt   int
	)

	switch initialValueType {
	case "date":
		initialValueDate, _ = time.Parse("2006-01-02", initialValue.(string))
	default:
		return nil, fmt.Errorf("no logic for value_type: %s (%v)", initialValueType, initialValue)
	}

	switch secondValueType {
	case "date":
		secondValueDate, _ = time.Parse("2006-01-02", secondValue.(string))
	case "int":
		secondValueInt = int(secondValue.(float64))
	default:
		return nil, fmt.Errorf("no logic for value_type: %s", secondValueType)
	}

	switch operation {
	case "subtract":
		return initialValueDate.AddDate(-secondValueInt, 0, 0).Format("2006-02-01"), nil
	case "cmp":
		if secondValueDate.Before(initialValueDate) {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		return nil, fmt.Errorf("no logic for operation: %s", operation)
	}
}

func waitingForValue(key string, selfMap *sync.Map) interface{} {
	if value, ok := selfMap.Load(key); ok {
		return value
	}
	return nil
}
