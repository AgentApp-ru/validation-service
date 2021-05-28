package main

import (
	"encoding/json"
	"validation_service/internal/apiserver/views"
)

func main() {
	a := []byte(`{
		"type": "date",
		"value": "1930-01-01"
	}`)

	var j views.DateMinMaxPattern
	var v views.DateDateValue

	json.Unmarshal(a, &j)
	println(j.PatternType)
	json.Unmarshal(j.Value, &v)
	println(v)
}
