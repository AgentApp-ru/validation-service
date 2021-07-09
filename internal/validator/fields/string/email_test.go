package string

// const emailPatterns = `
// [
//   {
//     "name": "local-part",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       }
//     ]
//   },
//   {
//     "name": "local-part",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       }
//     ]
//   },
//   {
//     "name": "local-part",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       }
//     ]
//   },
//   {
//     "name": "local-part",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9_+-]",
//         "max": 68,
//         "min": 1
//       }
//     ]
//   },
//   {
//     "name": "divider",
//     "patterns": [
//       {
//         "chars": "[@]",
//         "max": 1,
//         "min": 1
//       }
//     ]
//   },
//   {
//     "name": "domain",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       }
//     ]
//   },
//   {
//     "name": "domain",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       }
//     ]
//   },
//   {
//     "name": "domain",
//     "patterns": [
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       },
//       {
//         "chars": "[\\.]",
//         "max": 1,
//         "min": 1
//       },
//       {
//         "chars": "[a-zA-Z0-9-]",
//         "max": 18,
//         "min": 2
//       }
//     ]
//   }
// ]
// `

// func TestValidateValidEmail(t *testing.T) {
// 	field := "seogwipo_helipor_244490@instomat.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if !isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be valid email: '%s'", field)
// 	}
// }

// func TestValidateValidEmailWithDots(t *testing.T) {
// 	field := "d.oqwi@op-dw.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if !isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be valid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailWithInValidCharacters(t *testing.T) {
// 	field := "doqwi$%@o$%pdw.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailTooLongBeforePaw(t *testing.T) {
// 	field := "doqwisadsadsadsadwqdwqeqwe23ewqd123@opdw.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailTooShortBeforePaw(t *testing.T) {
// 	field := "a@opdw.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailTooLongAfterPaw(t *testing.T) {
// 	field := "aasd@oasdsadasdsadasdasdadadaasdasdpdw.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailTooShortAfterPaw(t *testing.T) {
// 	field := "sadasda@o.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailWithoutDot(t *testing.T) {
// 	field := "sadasda@asdoasdru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateInValidEmailWithoutPaw(t *testing.T) {
// 	field := "sadasdao.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateValidEmailOnlyNumbers(t *testing.T) {
// 	field := "123123@mail.ru"

// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }

// func TestValidateValidEmailOnlySpecialCharacters(t *testing.T) {
// 	field := "-.+-@mail.ru"
// 	var patterns []*StringPattern
// 	json.Unmarshal([]byte(json.RawMessage(emailPatterns)), &patterns)

// 	if isValidatedWithGroups(field, patterns, false) {
// 		t.Errorf("should be invalid email: '%s'", field)
// 	}
// }
