package string

// func TestValidateValidString(t *testing.T) {
// 	field := unicodeString("русский")
// 	min := 3

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: &min,
// 			Max:    10,
// 		},
// 	}

// 	if leftBody, ok := validateStringWithPatterns(field, patterns); !ok || len(leftBody) != 0 {
// 		pattern := patterns[0]
// 		t.Errorf("should be valid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateTooLongString(t *testing.T) {
// 	field := unicodeString("русский")
// 	min := 3

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: &min,
// 			Max:    4,
// 		},
// 	}

// 	if leftBody, ok := validateStringWithPatterns(field, patterns); ok && len(leftBody) == 0 {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateTooShortString(t *testing.T) {
// 	field := unicodeString("русский")
// 	min := 10

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: &min,
// 			Max:    14,
// 		},
// 	}

// 	if _, ok := validateStringWithPatterns(field, patterns); ok {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateValidStringWithExactSize(t *testing.T) {
// 	field := unicodeString("русский")

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: nil,
// 			Max:    7,
// 		},
// 	}

// 	if leftBody, ok := validateStringWithPatterns(field, patterns); !ok || len(leftBody) != 0 {
// 		pattern := patterns[0]
// 		t.Errorf("should be valid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateInValidStringWithExactSize(t *testing.T) {
// 	field := unicodeString("русск")

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: nil,
// 			Max:    7,
// 		},
// 	}

// 	if _, ok := validateStringWithPatterns(field, patterns); ok {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateInValidStringWithInvalidLetters(t *testing.T) {
// 	field := unicodeString("english")
// 	min := 0

// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[а-яА-ЯёЁ\\s-'`]",
// 			MinPtr: &min,
// 			Max:    10,
// 		},
// 	}

// 	if leftBody, ok := validateStringWithPatterns(field, patterns); ok && len(leftBody) == 0 {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(field))
// 	}
// }

// func TestValidateValidPhone(t *testing.T) {
// 	field := "79999999999"
// 	min := 11
// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[0-9]",
// 			MinPtr: &min,
// 			Max:    18,
// 		},
// 	}
// 	if _, ok := validateStringWithPatterns(unicodeString(field), patterns); !ok {
// 		pattern := patterns[0]
// 		t.Errorf("should be valid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }

// func TestPrepareDirtyPhone(t *testing.T) {
// 	field := "+7 (999) 123-45-67"
// 	transformers := `
// 	{
// 		"remove_chars": {
// 		  "chars": "[-+()\\s]"
// 		}
// 	  }
// 	`

// 	cleanField := prepare(field, json.RawMessage(transformers))

// 	if cleanField != "79991234567" {
// 		t.Errorf("should be clean field: '%s' for '%s'", cleanField, field)
// 	}
// }

// func TestPrepareDirtyPhoneWithoutTransformers(t *testing.T) {
// 	field := "+7 (999) 123-45-67"
// 	transformers := ``

// 	cleanField := prepare(field, json.RawMessage(transformers))

// 	if cleanField != "+7 (999) 123-45-67" {
// 		t.Errorf("should be clean field: '%s' for '%s'", cleanField, field)
// 	}
// }

// func TestValidateInValidPhoneInValidCharacters(t *testing.T) {
// 	field := "799999asd999-9+9"
// 	min := 11
// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[0-9]",
// 			MinPtr: &min,
// 			Max:    18,
// 		},
// 	}
// 	if _, ok := validateStringWithPatterns(unicodeString(field), patterns); ok {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }

// func TestValidateInValidPhoneTooSmall(t *testing.T) {
// 	field := "799999"
// 	min := 11
// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[0-9]",
// 			MinPtr: &min,
// 			Max:    18,
// 		},
// 	}
// 	if _, ok := validateStringWithPatterns(unicodeString(field), patterns); ok {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }

// func TestValidateInValidPhoneTooLarge(t *testing.T) {
// 	field := "7999999999999999999"
// 	min := 11
// 	patterns := []*Pattern{
// 		{
// 			Chars:  "[0-9]",
// 			MinPtr: &min,
// 			Max:    18,
// 		},
// 	}
// 	if leftBody, _ := validateStringWithPatterns(unicodeString(field), patterns); len(leftBody) == 0 {
// 		pattern := patterns[0]
// 		t.Errorf("should be invalid regexp: '%s' for '%s'", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), field)
// 	}
// }
