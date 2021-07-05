package string

import (
	"encoding/json"
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


func TestValidateValidEmail(t *testing.T) {
	field := "doqwi@opdw.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}

	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestAnotherValidateValidEmail(t *testing.T) {
	field := "testmail@mail.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}

	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailWithInValidCharacters(t *testing.T) {
	field := "doqwi$%@o$%pdw.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailTooLongBeforePaw(t *testing.T) {
	field := "doqwisadsadsadsadwqdwqeqwe23ewqd123@opdw.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailTooShortBeforePaw(t *testing.T) {
	field := "a@opdw.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailTooLongAfterPaw(t *testing.T) {
	field := "aasd@oasdsadasdsadasdasdadadaasdasdpdw.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailTooShortAfterPaw(t *testing.T) {
	field := "sadasda@o.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailWithoutDot(t *testing.T) {
	field := "sadasda@asdoasdru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidEmailWithoutPaw(t *testing.T) {
	field := "sadasdao.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateValidEmailOnlyNumbers(t *testing.T) {
	field := "123123@mail.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateValidEmailOnlySpecialCharacters(t *testing.T) {
	field := "-.+-@mail.ru"
	min_first := 2
	min_second := 1
	patterns := []*Pattern{
		{
			Chars:  "[a-zA-Z0-9_.+-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[@]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-]",
			MinPtr: &min_first,
			Max:    18,
		},
		{
			Chars:  "[.]",
			MinPtr: &min_second,
			Max:    1,
		},
		{
			Chars:  "[a-zA-Z0-9-.]",
			MinPtr: &min_first,
			Max:    18,
		},
	}
	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateValidPhone(t *testing.T) {
	field := "79999999999"
	min := 11
	patterns := []*Pattern{
		{
			Chars:  "[0-9]",
			MinPtr: &min,
			Max:    18,
		},
	}
	if !validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestPrepareDirtyPhone(t *testing.T) {
	field := "+7 (999) 123-45-67"
	transformers := `
	{
		"remove_chars": {
		  "chars": "[-+()\\s]"
		}
	  }
	`

	cleanField := prepare(field, json.RawMessage(transformers))

	if cleanField != "79991234567" {
		t.Errorf("should be clean field: '%s' for '%s'", cleanField, field)
	}
}

func TestPrepareDirtyPhoneWithoutTransformers(t *testing.T) {
	field := "+7 (999) 123-45-67"
	transformers := ``

	cleanField := prepare(field, json.RawMessage(transformers))

	if cleanField != "+7 (999) 123-45-67" {
		t.Errorf("should be clean field: '%s' for '%s'", cleanField, field)
	}
}

func TestValidateInValidPhoneInValidCharacters(t *testing.T) {
	field := "799999asd999-9+9"
	min := 11
	patterns := []*Pattern{
		{
			Chars:  "[0-9]",
			MinPtr: &min,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidPhoneTooSmall(t *testing.T) {
	field := "799999"
	min := 11
	patterns := []*Pattern{
		{
			Chars:  "[0-9]",
			MinPtr: &min,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}

func TestValidateInValidPhoneTooLarge(t *testing.T) {
	field := "7999999999999999999"
	min := 11
	patterns := []*Pattern{
		{
			Chars:  "[0-9]",
			MinPtr: &min,
			Max:    18,
		},
	}
	if validateStringWithPatterns(field, patterns) {
		pattern := patterns[0]
		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
	}
}
