package main

import (
	"fmt"
	"regexp"
)

func main() {
	// pattern := "[а-яА-ЯёЁ -`]"
	// stringToCheck := []byte(string([]rune("Римский-'`Корсаков")))

	pattern := "[а-яА-ЯёЁ '`-]"
	re := regexp.MustCompile(fmt.Sprintf("%s{%d,%d}", pattern, 5, 5))
	stringToCheck := []byte(string([]rune("а-'`а")))

	println(re.Match(stringToCheck))

	// matched, _ := regexp.Match(
	// 	fmt.Sprintf("%s{%d,%d}", pattern, 5, 5), stringToCheck,
	// )
	//
	// println(matched)
}
