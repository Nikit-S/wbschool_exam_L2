package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
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

type tokenArr struct {
	main [][][]byte
	sort [][][]byte
	free [][][]byte
	k    int
}

func (s tokenArr) Len() int       { return len(s.sort) }
func (s tokenArr) Swap(i, j int)  { s.sort[i], s.sort[j] = s.sort[j], s.sort[i] }
func (s tokenArr) Truncate(n int) { s.sort = s.sort[:n] }

type ByLetterDesc struct{ *tokenArr }

func (o ByLetterDesc) Less(i, j int) bool {
	return bytes.Compare(o.sort[i][o.k], o.sort[j][o.k]) == -1
}

type ByLetterAsc struct{ *tokenArr }

func (o ByLetterAsc) Less(i, j int) bool {
	return bytes.Compare(o.sort[i][o.k], o.sort[j][o.k]) == 1
}

func (s *tokenArr) getFreeSort(k int) {
	s.k = k
	index := 0
	for i := range s.main {
		if s.k >= len(s.main[i]) {
			s.main[i], s.main[index] = s.main[index], s.main[i]
			index++
		}
	}
	s.free = s.main[:index]
	s.sort = s.main[index:]
}

func Unique(s *tokenArr) {
	for i := range s.sort[:len(s.sort)-2] {
		if bytes.Compare(s.sort[i][s.k], s.sort[i+1][s.k]) == 0 {
			fmt.Println(s.sort[i][s.k], s.sort[i+1][s.k])
			s.sort = append(s.sort[:i], s.sort[i+1:]...)
			i++
		}
	}
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
	fmt.Println()

	file, err := os.Open(args[0])
	defer file.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	mainObj := tokenArr{}

	for scanner.Scan() {
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return
		}
		mainObj.main = append(mainObj.main, bytes.Split(lines, []byte(" ")))
	}

	mainObj.getFreeSort(*colNum - 1)

	if *revSort {
		sort.Sort(ByLetterAsc{&mainObj})
		mainObj.main = append(mainObj.sort, mainObj.free...)
	} else {
		sort.Sort(ByLetterDesc{&mainObj})
	}

	if *dupSort {
		fmt.Println("Unique")
		Unique(&mainObj)
	}

	for _, v := range mainObj.main {
		for _, k := range v {
			fmt.Print(k)
		}
		fmt.Println("")
	}

}
