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

// подгонка под интефейс sort
func (s stringArr) Less(i, j int) bool { return s[i] < s[j] }
func (s stringArr) Len() int           { return len(s) }
func (s stringArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type anagrams map[string][]string

//тут мы храним жанные в формате
// [пятак] = {[п]=1 [я]=1 [т]=1 [а]=1 [к]=1}
// так мы можем сократить количество памяти в случаях если унас в исходном слове
// много повтороений
type checker map[string]map[rune]int

type task struct {
	anagram map[string][]string
	checker checker
}

//приведение «пятак» к виду {[п]=1 [я]=1 [т]=1 [а]=1 [к]=1}
func makeRuneMap(str string) map[rune]int {
	retMap := make(map[rune]int)
	for _, letter := range str {
		retMap[letter] += 1
	}
	return retMap
}

//сложнаааааааааааааааааааааааааааааааааааа
func (t *task) IsAnagram(str string) (string, bool) {
	tempMap := makeRuneMap(str)

	//бежим по всем словам которые уже имеются в списке
	for key := range t.checker {
		b := false

		//бежим по всем буквам в каждом слове с которым встречались
		for letter := range t.checker[key] {
			v, ok := tempMap[letter]
			//если такая буква есть в проверямеом слове
			if ok {
				//нужно проверить сколько раз она встречается
				if v != t.checker[key][letter] {
					//если количество не совпадает
					b = false
					break
				}
				b = true
			} else {
				//такой буквы вообще нет в словаре
				b = false
				break
			}
		}
		// если нашли то отправялем обратно успех и то слово из которого мы
		// смогли составить анаграмму
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
		//игнорим регистр
		word = strings.ToLower(word)

		//проверяем есть лои это слово в списке на проверку
		_, baseWord := mainObj.checker[word]

		// если его нет значит мы с таким еще не встречались, начинаем работу
		if !baseWord {
			str, inCheck := mainObj.IsAnagram(word)
			if inCheck {
				//если слово возможно состваить добавляем его в соответсвующий массив
				mainObj.anagram[str] = append(mainObj.anagram[str], word)
			} else {
				// если слово составить нельзя значит оно уникальное и
				// добавляется в обе группы
				mainObj.checker[word] = makeRuneMap(word)
				mainObj.anagram[word] = append(mainObj.anagram[word], word)
			}
		}
	}

	//удаляем дубликаты, массивы в которых по одному слову, и сортируем
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
