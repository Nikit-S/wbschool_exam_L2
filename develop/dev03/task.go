package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func remsave(arr [][][]byte, i int) [][][]byte {
	return append((arr)[0:i], (arr)[i+1:]...)
}

func Uniq(a [][][]byte, k int) [][][]byte {
	fmt.Println("Uniq")

	i := 0
	for i < len(a)-2 {
		eq := true
		if k <= len(a[i])-1 && len(a[i]) == len(a[i+1]) {
			fmt.Println("Checking", a[i], a[i+1])
			for l := range a[i+1][k:] {
				if !bytes.Equal(a[i][k+l], a[i+1][k+l]) {
					eq = false
					break
				}
			}
			if eq {
				fmt.Println("removing ", i+1)
				a = remsave(a, i+1)
				i--
			}
		} else if len(a[i]) == len(a[i+1]) {
			fmt.Println("removing ", i+1)
			a = remsave(a, i+1)
			i--
		}
		i++
	}
	return a
}

func Less(a, b []byte) bool {
	for i, v := range a {
		if v == b[i] {
			if i == len(a)-1 || i == len(b)-1 {
				return len(a) < len(b)
			}
			continue
		} else if v < b[i] {
			return true
		} else if v > b[i] {
			return false
		}
	}
	return false
}

func Greater(a, b []byte) bool {
	for i, v := range a {
		if v == b[i] {
			if i == len(a)-1 || i == len(b)-1 {
				return len(a) > len(b)
			}
			continue
		} else if v > b[i] {
			return true
		} else if v < b[i] {
			return false
		}
	}
	return false
}

func quicksort(a [][][]byte, f func(a1, a2 []byte) bool, k int) [][][]byte {
	rand.Seed(time.Now().UnixNano())

	//сразу делаем возврат если в слайсе один элемент == все отсортировано
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	//ставим поворотный элемент в конец
	a[pivot], a[right] = a[right], a[pivot]

	//двигаемся и по ходу при необходимости меняем элементы меньшие с последним
	//запомнившимся элементом который больше
	for i := range a {
		if f(a[i][k], a[right][k]) {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	//в конце меняем поворотный на полседний что больше и получаем слайс в котором
	// все меньшие элементы сдева а все большие слева, с ними дальше и работаем
	/*for _, v := range a {
		fmt.Println(v)
	}*/
	fmt.Println(left, right)
	a[left], a[right] = a[right], a[left]
	quicksort(a[:left], f, k)
	quicksort(a[left+1:], f, k)

	return a
}

func main() {
	colNum := flag.Int("k", 1, "указание колонки для сортировки")
	numSort := flag.Bool("n", false, "сортировать по числовому значению")
	revSort := flag.Bool("r", false, "сортировать в обратном порядке")
	dupSort := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	args := flag.Args()
	fmt.Println("colNum: ", *colNum)
	fmt.Println("numSort: ", *numSort)
	fmt.Println("revSort: ", *revSort)
	fmt.Println("dupSort: ", *dupSort)
	fmt.Println("args: ", args)

	file, err := os.Open(args[0])
	defer file.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	lineWordsArray := [][][]byte{}

	for scanner.Scan() {
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return
		}
		lineWordsArray = append(lineWordsArray, bytes.Fields(lines))
	}

	//fmt.Println(v, len(lineWordsArray))
	for _, v := range lineWordsArray {
		fmt.Println(v)
	}

	if !*revSort {
		index := 0
		for i := range lineWordsArray {
			if *colNum-1 >= len(lineWordsArray[i]) {
				lineWordsArray[i], lineWordsArray[index] = lineWordsArray[index], lineWordsArray[i]
				index++
			}
		}
		lineWordsArray = append(lineWordsArray[:index], quicksort(lineWordsArray[index:], Less, *colNum-1)...)

	} else {
		index := len(lineWordsArray) - 1
		i := index
		for i >= 0 {
			if *colNum-1 >= len(lineWordsArray[i]) {
				lineWordsArray[i], lineWordsArray[index] = lineWordsArray[index], lineWordsArray[i]
				index--
			}
			i--
		}
		index++
		lineWordsArray = append(quicksort(lineWordsArray[0:index], Greater, *colNum-1), lineWordsArray[index:]...)
	}
	fmt.Println("-----------")
	for _, v := range lineWordsArray {
		fmt.Println(v)
	}
	if *dupSort {
		lineWordsArray = Uniq(lineWordsArray, *colNum-1)
	}

	fmt.Println("-----------")
	for _, v := range lineWordsArray {
		fmt.Println(v)
	}
}
