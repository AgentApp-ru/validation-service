package string

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/log"
)

type Pattern struct {
	Chars   string  `json:"chars"`
	MinPtr  *int    `json:"min"`
	Min     int     `json:"-"`
	Max     int     `json:"max"`
	Extract *string `json:"extract"`
}

type StringPattern struct {
	Name               string     `json:"name"`
	Allow_white_spaces bool       `json:"allow_white_spaces"`
	Patterns           []*Pattern `json:"patterns"`
}

func Validate(field interface{}, fieldValidator *fields.FieldValidator) interface{} {
	var (
		stringPatterns []*StringPattern
		ok             bool
		strField       string
	)

	if strField, ok = field.(string); !ok {
		log.Logger.Error("type conversion failed")
		return nil
	}

	if err := json.Unmarshal([]byte(fieldValidator.Patterns), &stringPatterns); err != nil {
		return nil
	}

	for _, stringPattern := range stringPatterns {
		if validateStringWithPatterns(strField, stringPattern.Patterns) {
			return strField
		}
	}
	return nil
}

func validateStringWithPatterns(field string, patterns []*Pattern) bool {
	leftBody := []rune(field)
	for _, pattern := range patterns {
		if pattern.MinPtr == nil {
			pattern.Min = pattern.Max
		} else {
			pattern.Min = *pattern.MinPtr
		}
		if pattern.Extract != nil {
			cleanString, err := deleteSpareCharacters(leftBody, *pattern.Extract)
			if err {
				log.Logger.Error("cleaning string failed")
				return false
			}
			leftBody = cleanString
		}
		if len(leftBody) < pattern.Min {
			return false
		}
		lenToCheck := int(math.Min(float64(len(leftBody)), float64(pattern.Max)))
		stringToCheck := []byte(string(leftBody[:lenToCheck]))
		minDimensionToCheck := int(
			math.Min(
				math.Min(float64(len(leftBody)), float64(pattern.Min)),
				float64(pattern.Max),
			),
		)
		matched, err := regexp.Match(
			fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max), stringToCheck,
		)
		if !matched || err != nil {
			return false
		}
		// После проверки строк на совпадение, отрезаем длину совпавшей части
		searching, err := regexp.Compile(
			fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max),
		)
		if err != nil {
			return false
		}
		cutting := searching.FindString(string(stringToCheck))
		cuttingLen := len([]rune(cutting))
		leftBody = leftBody[cuttingLen:]
	}

	return len(leftBody) == 0
}

func deleteSpareCharacters(field []rune, pattern string) ([]rune, bool) {
	searching, err := regexp.Compile(
		pattern,
	)
	if err != nil {
		return nil, true
	}
	cuttingField := string(field)
	cuttingField = searching.ReplaceAllString(cuttingField, "")
	return []rune(cuttingField), false
}
