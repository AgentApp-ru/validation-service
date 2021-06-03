package string

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/log"
)

type (
	unicodeString []rune

	Pattern struct {
		Chars  string `json:"chars"`
		MinPtr *int   `json:"min"`
		Min    int    `json:"-"`
		Max    int    `json:"max"`
	}

	StringPattern struct {
		Name     string     `json:"name"`
		Patterns []*Pattern `json:"patterns"`
	}

	CharsToRemove struct {
		Chars string `json:"chars"`
	}

	Transformers struct {
		CharsToRemove *CharsToRemove `json:"remove_chars"`
	}
)

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

	var preparedField string
	if fieldValidator.Transformers != nil {
		preparedField = prepare(strField, *fieldValidator.Transformers)
	} else {
		preparedField = strField
	}

	if err := json.Unmarshal([]byte(fieldValidator.Patterns), &stringPatterns); err != nil {
		return nil
	}

	if isValidatedWithGroups(preparedField, stringPatterns, fieldValidator.Allow_white_spaces) {
		return true
	}

	return nil
}

func prepare(field string, rawTransformers json.RawMessage) string {
	var (
		transformers *Transformers
		err          error
	)

	if err = json.Unmarshal([]byte(rawTransformers), &transformers); err != nil {
		return field
	}

	if transformers != nil && transformers.CharsToRemove != nil && transformers.CharsToRemove.Chars != "" {
		return deleteSpareCharacters(field, transformers.CharsToRemove.Chars)
	}
	return field
}

func in(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func isValidatedWithGroups(preparedField string, stringPatterns []*StringPattern, allowWhiteSpaces bool) bool {
	groups := []string{}
	for _, stringPattern := range stringPatterns {
		if !in(groups, stringPattern.Name) {
			groups = append(groups, stringPattern.Name)
		}
	}

	leftBodies := make(map[int][]unicodeString)
	leftBodies[0] = []unicodeString{
		unicodeString(preparedField),
	}

	for i, group := range groups {
		leftBodies[i+1] = []unicodeString{}
		for _, stringPattern := range stringPatterns {
			if stringPattern.Name != group {
				continue
			}

			for _, stringToCheck := range leftBodies[i] {
				leftBody, ok := validateStringWithPatterns(stringToCheck, stringPattern.Patterns)
				if ok {
					leftBodies[i+1] = append(leftBodies[i], leftBody)
					if len(leftBody) == 0 {
						if i+1 == len(groups) {
							return true
						}
					} else if allowWhiteSpaces && string(leftBody[0]) == " " {
						leftBodies[i+1] = append(leftBodies[i], leftBody[1:])
					}
				}
			}
		}
	}

	return false
}

func validateStringWithPatterns(leftBody unicodeString, patterns []*Pattern) (unicodeString, bool) {
	for _, pattern := range patterns {
		if pattern.MinPtr == nil {
			pattern.Min = pattern.Max
		} else {
			pattern.Min = *pattern.MinPtr
		}
		// println("!!", pattern.MinPtr == nil, pattern.Min)

		if len(leftBody) < pattern.Min {
			return nil, false
		}
		lenToCheck := int(math.Min(float64(len(leftBody)), float64(pattern.Max)))
		stringToCheck := []byte(string(leftBody[:lenToCheck]))
		minDimensionToCheck := int(
			math.Min(
				math.Min(float64(len(leftBody)), float64(pattern.Min)),
				float64(pattern.Max),
			),
		)
		// println(fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max), string(stringToCheck))
		matched, err := regexp.Match(
			fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max), stringToCheck,
		)
		if !matched || err != nil {
			return nil, false
		}
		// После проверки строк на совпадение, отрезаем длину совпавшей части
		searching, err := regexp.Compile(
			fmt.Sprintf("%s{%d,%d}", pattern.Chars, minDimensionToCheck, pattern.Max),
		)
		if err != nil {
			return nil, false
		}
		cutting := searching.FindString(string(stringToCheck))
		cuttingLen := len([]rune(cutting))
		leftBody = leftBody[cuttingLen:]
	}

	return leftBody, true
}

func deleteSpareCharacters(field, pattern string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return field
	}

	return reg.ReplaceAllString(field, "")
}
