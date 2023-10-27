package main

import "log"

func main() {
	log.Println("Hello World")
	log.Println(testAndFuzzMe(1, 2))
}

func testAndFuzzMe(x, y int) int {
	return x + y
}
