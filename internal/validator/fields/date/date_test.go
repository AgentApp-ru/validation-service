package date

// import (
// 	"encoding/json"
// 	"testing"
// 	"time"
// )

// func TestValidateInvalidMinDateDate(t *testing.T) {
// 	field := "2015-01-01"
// 	minField := "2018-01-01"

// 	fieldDate, _ := time.Parse("2006-01-02", field)
// 	rawMinValue := json.RawMessage(minField)

// 	if validateMinDate(fieldDate, rawMinValue) {
// 		t.Errorf("%s should be less than %s", minField, field)
// 	}
// }

// func TestValidateValidMinDateDate(t *testing.T) {
// 	field := "2020-01-01"
// 	minField := `"2018-01-01"`

// 	fieldDate, _ := time.Parse("2006-01-02", field)
// 	rawMinValue := json.RawMessage(minField)

// 	if !validateMinDate(fieldDate, rawMinValue) {
// 		t.Errorf("%s should be less than %s", minField, field)
// 	}
// }
