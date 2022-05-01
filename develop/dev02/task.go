package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
func countNu(st []rune, i int) int {
	for i < len(st) && unicode.IsDigit(st[i]) {
		i++
	}
	return i
}

func check(st []rune) error {
	if unicode.IsDigit(st[0]) {
		return fmt.Errorf("Wrong format")
	}
	if unicode.IsDigit(st[len(st)-1]) && len(st) > 1 && unicode.IsDigit(st[0]) {
		return fmt.Errorf("Wrong format")
	}
	return nil
}

// выясняем следующий не бекслеш и является ли числом
func escChar(st []rune, i int) int {
	if i < len(st)-1 && unicode.IsDigit(st[i+1]) || st[i+1] == '\\' {
		i++
	}
	return i
}

func unPack(st []rune) ([]rune, error) {
	var sb strings.Builder
	var err error
	if len(st) == 0 {
		return []rune(""), nil
	}
	if err = check(st); err != nil {
		return nil, err
	}

	for i := 0; i < len(st); i++ {
		// бежим по строке до бекслеша
		if st[i] == '\\' {
			i = escChar(st, i)
			sb.WriteRune(st[i])
			continue
		}
		//если натыкаемся на число, записываем следующий за ним символ нужное
		//количество раз и двигаем индекс на длину этого числа
		if st[i] >= '1' && st[i] < '9' {
			t := countNu(st, i)
			nu, _ := strconv.Atoi(string(st[i:t]))

			for nu > 1 {
				sb.WriteRune(st[i-1])
				nu--
			}
			i = t - 1
			continue
		}
		sb.WriteRune(st[i])
	}
	return []rune(sb.String()), nil
}

func main() {

}
