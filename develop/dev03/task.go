package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"unicode"

	"github.com/mpvl/unique"
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
	main [][]byte
	sort [][]byte
	k    int
}

func (s tokenArr) Len() int { return len(s.sort) }
func (s tokenArr) Swap(i, j int) {
	s.main[i], s.main[j] = s.main[j], s.main[i]
	s.sort[i], s.sort[j] = s.sort[j], s.sort[i]

}
func (s tokenArr) Truncate(n int) {
	s.main = s.main[:n-1]
	s.sort = s.sort[:n-1]
}

type ByLetterDesc struct{ *tokenArr }

func (o ByLetterDesc) Less(i, j int) bool {
	x, y := []byte{}, []byte{}
	if len(o.sort[i]) == 0 && len(o.sort[j]) == 0 {
		x = o.main[i]
		y = o.main[j]
	} else {
		x = o.sort[i]
		y = o.sort[j]
	}
	return bytes.Compare(x, y) == -1
}

type ByLetterAsc struct{ *tokenArr }

func (o ByLetterAsc) Less(i, j int) bool {
	x, y := []byte{}, []byte{}
	if len(o.sort[i]) == 0 && len(o.sort[j]) == 0 {
		x = o.main[i]
		y = o.main[j]
	} else {
		x = o.sort[i]
		y = o.sort[j]
	}
	return bytes.Compare(x, y) == 1
}

func IsVisible(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsSpace(r)
}

type ByNumberDesc struct{ *tokenArr }

func (o ByNumberDesc) Less(i, j int) bool {
	x, y := []byte{}, []byte{}
	if len(o.sort[i]) == 0 && len(o.sort[j]) == 0 {
		x = o.main[i]
		y = o.main[j]
	} else {
		x = o.sort[i]
		y = o.sort[j]
	}
	a, err := strconv.Atoi(string(x))
	if err != nil {
		return bytes.Compare(x, y) == -1
	}
	b, err := strconv.Atoi(string(y))
	if err != nil {
		return bytes.Compare(x, y) == -1
	}
	return a < b
}

type ByNumberAsc struct{ *tokenArr }

func (o ByNumberAsc) Less(i, j int) bool {
	x, y := []byte{}, []byte{}
	if len(o.sort[i]) == 0 && len(o.sort[j]) == 0 {
		x = o.main[i]
		y = o.main[j]
	} else {
		x = o.sort[i]
		y = o.sort[j]
	}
	a, err := strconv.Atoi(string(x))
	if err != nil {
		return bytes.Compare(x, y) == 1
	}
	b, err := strconv.Atoi(string(y))
	if err != nil {
		return bytes.Compare(x, y) == 1
	}
	return a > b
}

func (s *tokenArr) getFreeSort(k int) {
	if k < 0 {
		os.Exit(1)
	}
	valid := false
	s.k = k
	for i := range s.main {
		temp := s.main[i]
		for j := 0; j < k; j++ {
			temp = bytes.TrimLeftFunc(temp, unicode.IsSpace)
			//fmt.Println(s.main[i])
			temp = bytes.TrimLeftFunc(temp, IsVisible)
			//fmt.Println(s.main[i])
		}
		if len(temp) > 0 {
			valid = true
		}
		s.sort = append(s.sort, temp)
		//fmt.Println(s.main[i])
		//fmt.Println(s.sort[i])
	}
	if !valid {
		s.k = 0
		for i := range s.main {
			s.sort[i] = s.main[i]
		}

	}
	//fmt.Println(valid)

}

func Unique(s *tokenArr) {
	if len(s.main) <= 1 {
		return
	}
	i := 0
	for i < len(s.main)-1 {
		if bytes.Compare(s.main[i], s.main[i+1]) == 0 {
			//s.sort = append(s.main[:i], s.main[i+1:]...)
			s.main = append(s.main[:i], s.main[i+1:]...)
			//fmt.Println("--upmain--")
			//for _, v := range s.main {
			//	fmt.Println(string(v))
			//}
			//continue
		}
		i++
	}
}

func main() {
	colNum := flag.Int("k", 1, "указание колонки для сортировки")
	numSort := flag.Bool("n", false, "сортировать по числовому значению")
	revSort := flag.Bool("r", false, "сортировать в обратном порядке")
	dupSort := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	args := flag.Args()
	//fmt.Println("colNum: ", *colNum)
	//fmt.Println("numSort: ", *numSort)
	//fmt.Println("revSort: ", *revSort)
	//fmt.Println("dupSort: ", *dupSort)
	//fmt.Println("args: ", args)
	//fmt.Println()

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
		mainObj.main = append(mainObj.main, lines)
	}

	mainObj.getFreeSort(*colNum - 1)
	//for _, v := range mainObj.sort {
	//	fmt.Println(string(v))
	//}
	//fmt.Println("-------")
	//for _, v := range mainObj.main {
	//	fmt.Println(string(v))
	//}

	//for _, v := range mainObj.sort {
	//	fmt.Println(string(v))
	//}
	//fmt.Println("-----")
	var temObj unique.Interface
	if *revSort {
		if *numSort {
			temObj = &ByNumberAsc{&mainObj}
		} else {
			temObj = &ByLetterAsc{&mainObj}
		}
		//mainObj.main = append(mainObj.sort, mainObj.free...)
	} else {
		if *numSort {
			temObj = &ByNumberDesc{&mainObj}
		} else {
			temObj = &ByLetterDesc{&mainObj}
		}
	}
	sort.Sort(temObj)
	//for _, v := range mainObj.main {
	//	fmt.Println(string(v))
	//}
	if *dupSort {
		Unique(&mainObj)
		//fmt.Println("------main_unique-----")
	}
	//fmt.Println("---result---")
	for _, v := range mainObj.main {
		fmt.Println(string(v))
	}

}
