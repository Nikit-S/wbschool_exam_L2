package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	//fmt.Println(i)
	i[1] = "5"
	i = append(i, "6")
	//fmt.Println(i)
}
