package string

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"validation_service/internal/validator/fields"
)


type Pattern struct {
	Chars string `json:"chars"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
}

type StringPattern struct {
	Name               string     `json:"name"`
	Allow_white_spaces bool       `json:"allow_white_spaces"`
	Patterns           []*Pattern `json:"patterns"`
}

func Validate(field string, fieldValidator *fields.FieldValidator) bool {
	var (
		ok             bool
		stringPatterns []*StringPattern
	)

	err := json.Unmarshal([]byte(fieldValidator.Patterns), &stringPatterns)
	if err != nil {
		return false
	}

	for _, stringPattern := range stringPatterns {
		println(stringPattern.Name)
		ok = true
		leftBody := []rune(field)

		for _, pattern := range stringPattern.Patterns {
			// println(pattern.Chars)
			if pattern.Min == 0 {
				pattern.Min = pattern.Max
			}

			// check-size

			if len(leftBody) < pattern.Min {
				ok = false
				break
			}

			asd := int(math.Min(float64(len(leftBody)), float64(pattern.Max)))
			stringToCheck := []byte(string(leftBody[:asd]))
			// println("regexp ", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(leftBody), string(stringToCheck))

			matched, err := regexp.Match(fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), stringToCheck)
			if !matched || err != nil {
				ok = false
			}

			leftBody = leftBody[asd:]
		}

		if len(leftBody) > 0 {
			ok = false
		}

		if ok {
			return true
		}
	}

	return false
}
