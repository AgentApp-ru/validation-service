package date

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetDateDate(t *testing.T) {
	minField := `"2018-01-01"`
	rawMinValue := json.RawMessage(minField)

	expectedDate, _ := time.Parse("2006-01-02", "2018-01-01")

	if date, ok := getPlainDate(rawMinValue); !ok || !date.Equal(expectedDate) {
		t.Errorf("%v should be equal to %v", date, expectedDate)
	}
}
