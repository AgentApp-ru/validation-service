package main

import (
	"regexp"
)

func main() {
	a := "Е271ХМ178"
	b := []rune(a)
	c := b[:1]
	d := string(c)
	println(d)
	matched, err := regexp.Match("[АВЕКМНОРСТУХ]{1,1}", []byte(d))
	println(matched, err)
}
