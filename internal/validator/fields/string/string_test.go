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
