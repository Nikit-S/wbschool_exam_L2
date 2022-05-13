package main

import (
	"fmt"
	"os"
)

type goo interface {
}

type gooi struct{}

func (g *gooi) t() {}

func Foo() error {
	var err os.PathError
	return &err
}

func Goo() goo {
	var g gooi
	return g
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
	t := Goo()
	fmt.Println(t)
	fmt.Println(t == nil)
}

//package main
//
//import (
//	"fmt"
//	"os"
//)
//
//func Foo() error {
//	var err *os.PathError = nil
//	return err
//}
//
//func main() {
//	err := Foo()
//	fmt.Println(err)
//	fmt.Println(err == nil)
//}
