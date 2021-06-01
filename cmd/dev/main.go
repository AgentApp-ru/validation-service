package main

import (
	"regexp"
)

func main() {
	matched, err := regexp.Match("[a-zA-Z0-9_.+-]{2,18}", []byte("vf@b2bpolis.ru"))
	println(matched, err)
}
