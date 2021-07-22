package requirements

import (
	"encoding/json"
	"sync"
	"time"
)

type (
	isRequired bool
	requiredAlternatives []string
	dependsOn string

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

		return isWaitingForFieldSucceed(value, selfMap)
	case "depends_on":
		var value dependsOn
		err := json.Unmarshal(r.Required.Value, &value)
		if err != nil {
			// TODO: log error
			println("!!!", err.Error())
			return false
		}

		fields := make([]string, 0)

		if isWaitingForFieldSucceed(append(fields, string(value)), selfMap) {
			return isFieldPresent
		} else {
			return true
		}
	default:
		return false
	}
}

func isWaitingForFieldSucceed(fieldsToWait []string, scopeObjectMap *sync.Map) bool {
	end := time.Now().Add(2 * time.Second)

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
