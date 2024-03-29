package string

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
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

	Field struct {
		Name     string     `json:"name"`
		Patterns []*Pattern `json:"patterns"`
	}

	CharsToRemove struct {
		Chars string `json:"chars"`
	}

	Transformers struct {
		CharsToRemove *CharsToRemove `json:"remove_chars"`
	}

	Validator struct {
		fieldName                    string
		objectMap                    *sync.Map
		allFieldsMap                 *sync.Map
		errors                       chan string
		transformers                 *json.RawMessage
		patterns                     json.RawMessage
		allowWhiteSpaces             bool
		maxConsecutiveSimilarSymbols int
		maxDuplicatedSymbols         int
	}
)

func New() *Validator {
	return new(Validator)
}

func (sv *Validator) Init(
	fieldName string,
	objectMap,
	allFieldsMap *sync.Map,
	errors chan string,
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
	maxConsecutiveSimilarSymbols int,
	maxDuplicatedSymbols int,
) {
	sv.fieldName = fieldName
	sv.objectMap = objectMap
	sv.allFieldsMap = allFieldsMap
	sv.errors = errors
	sv.transformers = transformers
	sv.patterns = patterns
	sv.allowWhiteSpaces = allowWhiteSpaces
	sv.maxConsecutiveSimilarSymbols = maxConsecutiveSimilarSymbols
	sv.maxDuplicatedSymbols = maxDuplicatedSymbols
}

func (sv *Validator) Validate(field interface{}) bool {
	var (
		stringPatterns []*Field
		ok             bool
		strField       string
	)

	if strField, ok = field.(string); !ok {
		log.Logger.Errorf("type conversion failed: %s -> %v", sv.fieldName, field)
		return false
	}

	var preparedField string
	if sv.transformers != nil {
		preparedField = prepare(strField, *sv.transformers)
	} else {
		preparedField = strField
	}

	if preparedField == "" {
		return false
	}

	if sv.maxConsecutiveSimilarSymbols != 0 {
		if !isValidatedForSimilarSymbols(preparedField, sv.maxConsecutiveSimilarSymbols) {
			return false
		}
	}

	if sv.maxDuplicatedSymbols != 0 {
		if !isValidatedForDuplicatedSymbols(preparedField, sv.maxDuplicatedSymbols) {
			return false
		}
	}

	if err := json.Unmarshal([]byte(sv.patterns), &stringPatterns); err != nil {
		return false
	}

	if isValidatedWithGroups(preparedField, stringPatterns, sv.allowWhiteSpaces) {
		return true
	}

	return false
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

func isValidatedForDuplicatedSymbols(preparedField string, maxDuplicateSymbols int) bool {
	for _, char := range preparedField {
		if strings.Count(preparedField, string(char)) > maxDuplicateSymbols {
			return false
		}
	}
	return true
}

func isValidatedForSimilarSymbols(preparedField string, maxConsecutiveSimilarSymbols int) bool {
	var (
		currentChar rune
		count       = 0
	)

	for _, v := range preparedField {
		if v == currentChar {
			count++
		} else {
			count = 1
		}
		currentChar = v

		if count > maxConsecutiveSimilarSymbols {
			return false
		}
	}
	return true
}

func isValidatedWithGroups(preparedField string, stringPatterns []*Field, allowWhiteSpaces bool) bool {
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
