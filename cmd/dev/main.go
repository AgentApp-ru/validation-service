package main

func main() {
	var a interface{}

	a = 123

	println(a)
	println(a.(int))
}
