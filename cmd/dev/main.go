package main

import (
	"fmt"
	"time"
)

func main() {
	fieldDate, _ := time.Parse("2006-01-02", "2021-07-14")
	println(fmt.Sprintf("%v", fieldDate))

	// fieldDate = fieldDate.Truncate(24*time.Hour)
	// println(fmt.Sprintf("%v", fieldDate))

	t := time.Now()
	t = t.Truncate(24*time.Hour)

	println(fmt.Sprintf("%v", t))

	println(!fieldDate.After(t))
	println(!fieldDate.Before(t))
	println(fieldDate.Equal(t))
}
