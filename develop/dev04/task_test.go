package main

import (
	"testing"
)

type testCase struct {
	input  []string
	result map[string][]string
}

func Equal(a, b []string) bool {
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

var simple = []testCase{
	testCase{
		input: []string{"пятак", "листок", "пятка", "тяпка", "слиток", "столик", "столи", "абв", "бва", "бав", "ваб", "ваб", "абв"},
		result: map[string][]string{
			"абв":    []string{"абв", "бав", "бва", "ваб"},
			"листок": []string{"листок", "слиток", "столик"},
			"пятак":  []string{"пятак", "пятка", "тяпка"},
		},
	},
	testCase{
		input: []string{"фбв", "Гора", "агОр", "агрО", "аРгО", "ГаОр", "гарО", "гоар", "огар", "орга", "раго", "рога"},
		result: map[string][]string{
			"гора": []string{"агор", "агро", "арго", "гаор", "гаро", "гоар", "гора", "огар", "орга", "раго", "рога"},
		},
	},
}

func TestSimple(t *testing.T) {
	for i := range simple {
		obj := findAnagrams(simple[i].input)
		if len(obj) != len(simple[i].result) {
			t.Log("Different len", obj, simple[i].result)
			t.Fail()
		}
		for k := range obj {
			j, ok := simple[i].result[k]
			if !ok || !Equal(obj[k], j) {
				t.Log(obj[k], j)
				t.Fail()
			}
		}
	}
}
