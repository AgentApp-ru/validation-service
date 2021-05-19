package validator

// func TestCarValidationWithoutFileReturnsNil(t *testing.T) {
// 	var err error

// 	validationsPath, err = ioutil.TempDir("/tmp", "validations")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer os.Remove(validationsPath)

// 	data, err := GetValidation("car")
// 	if data != nil {
// 		t.Error("data is not nil")
// 	}

// 	if err.Error() != fmt.Sprintf("open %s/car.json: no such file or directory", validationsPath) {
// 		t.Error(err)
// 	}
// }

// func TestPersonValidationWithoutFileReturnsNil(t *testing.T) {
// 	var err error

// 	validationsPath, err = ioutil.TempDir("/tmp", "validations")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer os.Remove(validationsPath)

// 	data, err := GetValidation("person")
// 	if data != nil {
// 		t.Error("data is not nil")
// 	}

// 	if err.Error() != fmt.Sprintf("open %s/person.json: no such file or directory", validationsPath) {
// 		t.Error(err)
// 	}
// }

// func TestValidationWithFileReturnsJSON(t *testing.T) {
// 	var err error

// 	validationsPath, err = ioutil.TempDir("/tmp", "validations")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer os.Remove(validationsPath)
// 	fi, err := ioutil.TempFile(validationsPath, "car*.json")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	defer os.Remove(fi.Name())

// 	content := []byte(`{
// 	"field": "manufacturing_year",
// 	"type": "date",
// 	"max": [],
// 	"min": [
// 		"1930-01-01"
// 	]
// }`)
// 	fi.Write(content)

// 	pos := strings.LastIndex(path.Base(fi.Name()), ".")
// 	fileName := path.Base(fi.Name())[:pos]

// 	data, err := GetValidation(fileName)
// 	if err != nil || data == nil {
// 		t.Errorf("data is nil: %s", err)
// 	}

// 	dataStr, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		t.Fatal("error marshaling the result: ", err)
// 	}

// 	diffOpts := jsondiff.DefaultConsoleOptions()
// 	res, diff := jsondiff.Compare(content, []byte(dataStr), &diffOpts)

// 	if res != jsondiff.FullMatch {
// 		t.Errorf("the expected result is not equal to what we have: %s", diff)
// 	}
// }
