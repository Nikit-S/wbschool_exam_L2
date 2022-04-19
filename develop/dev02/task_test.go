package main

import (
	"fmt"
	"os"
	"testing"
)

type elem struct {
	test    []rune
	correct []rune
}

var testfunc = unPack

var sim_elems = []elem{
	elem{[]rune("Hello"), []rune("Hello")},
	elem{[]rune("Hel1lo"), []rune("Hello")},
	elem{[]rune("He4llo"), []rune("Heeeello")},
	elem{[]rune("He4l11lo"), []rune("Heeeellllllllllllo")},
	elem{[]rune("He4 l11lo"), []rune("Heeee llllllllllllo")},
	elem{[]rune("He4 4l11lo"), []rune("Heeee    llllllllllllo")},
	elem{[]rune(""), []rune("")},
	elem{[]rune{123, '4', 111}, []rune{123, 123, 123, 123, 111}},
	elem{[]rune{1000000, '4', 111}, []rune{1000000, 1000000, 1000000, 1000000, 111}},
}

func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestSimple(t *testing.T) {

	for _, e := range sim_elems {
		res, _ := testfunc(e.test)
		fmt.Println(res)
		fmt.Println(e.correct)
		if !Equal(res, e.correct) {
			os.Exit(1)
		}
	}
}

var sim_errors = []elem{
	elem{[]rune("45"), []rune("45")},
	elem{[]rune("5555"), []rune("5555")},
	elem{[]rune("4"), []rune("4")},
	elem{[]rune("e4a"), []rune("eeeea")},
	elem{[]rune("e4"), []rune("eeee")},
	elem{[]rune("4e"), []rune("eeee")},
}

func TestErrors(t *testing.T) {

	for _, e := range sim_errors {
		res, err := testfunc(e.test)
		if err != nil {
			fmt.Printf("|%s|: %s\n", string(e.test), err.Error())
			continue
		}
		if !Equal(res, e.correct) {
			fmt.Printf("t|%s|\nc|%s|\n", string(res), string(e.correct))
			os.Exit(1)
		}
	}
}

var sim_esc = []elem{
	elem{[]rune("He\\2o"), []rune("He2o")},
	elem{[]rune("He\\23o"), []rune("He222o")},
	elem{[]rune("He\\\\5o"), []rune("He\\\\\\\\\\o")},
	elem{[]rune("\\\\"), []rune("\\")},
}

func TestEscape(t *testing.T) {

	for _, e := range sim_esc {
		res, err := testfunc(e.test)
		if err != nil {
			fmt.Printf("|%s|: %s\n", string(e.test), err.Error())
			continue
		}
		fmt.Printf("t|%s|\nc|%s|\n", string(res), string(e.correct))
		if !Equal(res, e.correct) {
			os.Exit(1)
		}
	}
}
