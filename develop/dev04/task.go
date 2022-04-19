package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
"пятак", "пятка" и "тяпка" - принадлежат одному множеству,
"листок", "слиток" и "столик" - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
type stringArr []string

func (s stringArr) Less(i, j int) bool { return s[i] < s[j] }
func (s stringArr) Len() int           { return len(s) }
func (s stringArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type anagrams map[string][]string
type checker map[string]map[rune]int

type task struct {
	anagram map[string][]string
	checker checker
}

func makeRuneMap(str string) map[rune]int {
	retMap := make(map[rune]int)
	for _, letter := range str {
		retMap[letter] += 1
	}
	return retMap
}

func (t *task) IsAnagram(str string) (string, bool) {
	tempMap := makeRuneMap(str)
	for key := range t.checker {
		b := false
		for letter := range t.checker[key] {
			v, ok := tempMap[letter]
			if ok {
				if v != t.checker[key][letter] {
					b = false
					break
				}
				b = true
			} else {
				b = false
				break
			}
		}
		if b {
			return key, true
		}
	}
	return "", false
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func findAnagrams(dictionary []string) map[string][]string {
	mainObj := task{anagram: make(anagrams),
		checker: make(map[string]map[rune]int)}
	for _, word := range dictionary {
		word = strings.ToLower(word)
		_, anagramWord := mainObj.anagram[word]
		_, baseWord := mainObj.checker[word]
		if !anagramWord && !baseWord {
			str, inCheck := mainObj.IsAnagram(word)
			//fmt.Printf("Check for %v is %v, group [%v]\n", word, inCheck, str)
			if inCheck {
				//fmt.Printf("adding %v to [%v]\n", str, word)
				mainObj.anagram[str] = append(mainObj.anagram[str], word)
			} else {
				//fmt.Printf("adding %v to [%v]\n", word, word)
				mainObj.checker[word] = makeRuneMap(word)
				mainObj.anagram[word] = append(mainObj.anagram[word], word)
			}
		}
	}
	for key := range mainObj.anagram {
		if len(mainObj.anagram[key]) <= 1 {
			delete(mainObj.anagram, key)
			continue
		}
		mainObj.anagram[key] = unique(mainObj.anagram[key])
		sort.Sort(stringArr(mainObj.anagram[key]))
	}
	return map[string][]string(mainObj.anagram)
}

func main() {

	dictionary := []string{"пятак", "листок", "пятка", "тяпка", "слиток", "столик", "столи", "абв", "бва", "бав", "ваб", "ваб", "абв"}
	obj := findAnagrams(dictionary)
	fmt.Println(obj)

}
