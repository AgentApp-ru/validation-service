package string

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"validation_service/internal/validator/fields"
)

type Pattern struct {
	Chars  string `json:"chars"`
	MinPtr *int   `json:"min"`
	Min    int    `json:"-"`
	Max    int    `json:"max"`
}

type StringPattern struct {
	Name               string     `json:"name"`
	Allow_white_spaces bool       `json:"allow_white_spaces"`
	Patterns           []*Pattern `json:"patterns"`
}

func Validate(field string, fieldValidator *fields.FieldValidator) bool {
	var (
		stringPatterns []*StringPattern
	)

	err := json.Unmarshal([]byte(fieldValidator.Patterns), &stringPatterns)
	if err != nil {
		return false
	}

	for _, stringPattern := range stringPatterns {
		if validateStringWithPatterns(field, stringPattern.Patterns) {
			return true
		}
	}

	return false
}

func validateStringWithPatterns(field string, patterns []*Pattern) bool {
	leftBody := []rune(field)

	for _, pattern := range patterns {
		if pattern.MinPtr == nil {
			pattern.Min = pattern.Max
		} else {
			pattern.Min = *pattern.MinPtr
		}

		// check-size
		if len(leftBody) < pattern.Min {
			// println("no len")
			return false
		}

		lenToCheck := int(math.Min(float64(len(leftBody)), float64(pattern.Max)))
		stringToCheck := []byte(string(leftBody[:lenToCheck]))
		minDimensionToCheck := int(
			math.Min(
				math.Max(float64(len(leftBody)), float64(pattern.Min)),
				float64(pattern.Max),
			),
		)
		// println("regexp ", fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max), string(leftBody), string(stringToCheck))

		matched, err := regexp.Match(
			fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max), stringToCheck,
		)
		if !matched || err != nil {
			// println("no match")
			return false
		}

		leftBody = leftBody[lenToCheck:]
	}

	return len(leftBody) == 0
}
