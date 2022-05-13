package main

import (
	"fmt"
)

func test() (x int) {
	x = 1
	defer func() {
		x++
		fmt.Println("defer t", x)

	}()
	return
}

func anotherTest() int {
	var x int
	x = 1
	defer func() {
		x++
		fmt.Println("defer at", x)
	}()
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
