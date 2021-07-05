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
	Pattern struct {
		Chars   string  `json:"chars"`
		MinPtr  *int    `json:"min"`
		Min     int     `json:"-"`
		Max     int     `json:"max"`
		Extract *string `json:"extract"`
	}

	StringPattern struct {
		Name               string     `json:"name"`
		Allow_white_spaces bool       `json:"allow_white_spaces"`
		Patterns           []*Pattern `json:"patterns"`
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

	for _, stringPattern := range stringPatterns {
		if validateStringWithPatterns(preparedField, stringPattern.Patterns) {
			return preparedField
		}
	}
	return nil
}

func prepare(field string, rawTransformers json.RawMessage) string {
	var (
		transformers   *Transformers
		err            error
	)

	if err = json.Unmarshal([]byte(rawTransformers), &transformers); err != nil {
		return field
	}

	if transformers != nil && transformers.CharsToRemove != nil && transformers.CharsToRemove.Chars != "" {
		return deleteSpareCharacters(field, transformers.CharsToRemove.Chars)
	}
	return field
}

func validateStringWithPatterns(field string, patterns []*Pattern) bool {
	leftBody := []rune(field)
	for _, pattern := range patterns {
		if pattern.MinPtr == nil {
			pattern.Min = pattern.Max
		} else {
			pattern.Min = *pattern.MinPtr
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

func deleteSpareCharacters(field, pattern string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return field
	}

	return reg.ReplaceAllString(field, "")
}
