package main

import (
	"regexp"
)

func main() {
	matched, err := regexp.Match("[а-яА-ЯёЁ\\s-'`]{4,7}", []byte("sdfg"))
	println(matched, err)
}
