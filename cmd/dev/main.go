package main

import (
	"fmt"
	"regexp"
)

func main() {
	matched, err := regexp.Match(`[\dA-HJ-NPR-Z]{13,13}`, []byte(`TMBED45J2B320`))
	fmt.Println(matched, err)
}
