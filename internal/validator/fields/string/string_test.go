package string

import (
	"fmt"
	"testing"
)

func TestValidateValidString(t *testing.T) {
	field := "русский"
	min := 3

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: &min,
			Max:    10,
		},
	}

	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be valid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateTooLongString(t *testing.T) {
	field := "русский"
	min := 3

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: &min,
			Max:    4,
		},
	}

	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateTooShortString(t *testing.T) {
	field := "русский"
	min := 10

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: &min,
			Max:    14,
		},
	}

	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateValidStringWithExactSize(t *testing.T) {
	field := "русский"

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: nil,
			Max:    7,
		},
	}

	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidStringWithExactSize(t *testing.T) {
	field := "русский"

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: nil,
			Max:    7,
		},
	}

	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidStringWithInvalidLetters(t *testing.T) {
	field := "english"
	min := 0

	patterns := []*Pattern{
		{
			Chars:  "[а-яА-ЯёЁ\\s-'`]",
			MinPtr: &min,
			Max:    10,
		},
	}

	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

// func TestValidateValidEmail(t *testing.T) {
// 	field := "vf@b2bpolis.ru"
// 	min1 := 1
// 	min2 := 2

// 	patterns := []*Pattern{
// 		{
// 			Chars: "[a-zA-Z0-9_.+-]",
// 			Max: 18,
// 			MinPtr: &min2,
// 		  },
// 		  {
// 			Chars: "[@]",
// 			Max: 1,
// 			MinPtr: &min1,
// 		  },
// 		  {
// 			Chars: "[a-zA-Z0-9-]",
// 			Max: 18,
// 			MinPtr: &min2,
// 		  },
// 		  {
// 			Chars: "[.]",
// 			Max: 1,
// 			MinPtr: &min1,
// 		  },
// 		  {
// 			Chars: "[a-zA-Z0-9-.]",
// 			Max: 18,
// 			MinPtr: &min2,
// 		  },
// 	}

// 	if !validateStringWithPatterns(field, patterns) {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }

// func TestValidateValidVIN(t *testing.T) {
// 	field := "TMBED45J2B3209311"

// 	var stringPatterns []*StringPattern
// 	json.Unmarshal([]byte(validatorClass.FieldValidatorsMap["vin"].Patterns), &stringPatterns)

// 	patterns := stringPatterns[0].Patterns
// 	if !validateStringWithPatterns(field, patterns) {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }
