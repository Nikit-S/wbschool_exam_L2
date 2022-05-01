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
}

// далее идет наполнение инетфейсов под стандартный sort
func (s tokenArr) Len() int { return len(s.sort) }
func (s tokenArr) Swap(i, j int) {
	s.main[i], s.main[j] = s.main[j], s.main[i]
	s.sort[i], s.sort[j] = s.sort[j], s.sort[i]

}
func (s tokenArr) Truncate(n int) {
	s.main = s.main[:n-1]
	s.sort = s.sort[:n-1]
}

func (s tokenArr) chooseLess(i, j int) ([]byte, []byte) {
	x, y := []byte{}, []byte{}
	if len(s.sort[i]) == 0 && len(s.sort[j]) == 0 {
		x = s.main[i]
		y = s.main[j]
	} else {
		x = s.sort[i]
		y = s.sort[j]
	}
	return x, y
}

type ByLetterDesc struct{ *tokenArr }

func (o ByLetterDesc) Less(i, j int) bool {
	x, y := o.chooseLess(i, j)
	//мы принимаем решение о том как будем сортпировать
	//либо у нас есть К столбец и тогда сравнивается сортируемая часть
	//либо обрезанная
	return bytes.Compare(x, y) == -1
}

type ByLetterAsc struct{ *tokenArr }

func (o ByLetterAsc) Less(i, j int) bool { return !ByLetterDesc(o).Less(i, j) }

func IsVisible(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsSpace(r)
}

type ByNumberDesc struct{ *tokenArr }

func (o ByNumberDesc) Less(i, j int) bool {
	x, y := o.chooseLess(i, j)
	a, err1 := strconv.Atoi(string(x))
	b, err2 := strconv.Atoi(string(y))
	if err1 != nil && err2 != nil {
		return bytes.Compare(x, y) == -1
	}
	if err1 != nil && err2 == nil {
		return bytes.Compare(nil, y) == -1
	}
	if err1 == nil && err2 != nil {
		return bytes.Compare(x, nil) == -1
	}
	return a < b
}

type ByNumberAsc struct{ *tokenArr }

func (o ByNumberAsc) Less(i, j int) bool { return !ByNumberDesc(o).Less(i, j) }

func (s *tokenArr) getFreeSort(k int) {
	if k < 0 {
		os.Exit(1)
	}
	valid := false
	// если можем выделить К столбец то тогда он отправлятся в категорию sort
	//со всем предшествующими пробелами
	for i := range s.main {
		temp := s.main[i]
		for j := 0; j < k; j++ {
			temp = bytes.TrimLeftFunc(temp, unicode.IsSpace)
			temp = bytes.TrimLeftFunc(temp, IsVisible)
		}
		if len(temp) > 0 {
			valid = true
		}
		s.sort = append(s.sort, temp)
	}
	if !valid {
		for i := range s.main {
			s.sort[i] = s.main[i]
		}
	}

}

func Unique(s *tokenArr) {
	if len(s.main) <= 1 {
		return
	}
	i := 0
	for i < len(s.main)-1 {

		if bytes.Compare(s.main[i], s.main[i+1]) == 0 {
			s.main = append(s.main[:i], s.main[i+1:]...)
			continue
		}
		i++
	}
}

func main() {
	//блок инициализации
	colNum := flag.Int("k", 1, "указание колонки для сортировки")
	numSort := flag.Bool("n", false, "сортировать по числовому значению")
	revSort := flag.Bool("r", false, "сортировать в обратном порядке")
	dupSort := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		return
	}
	file, err := os.Open(args[0])
	defer file.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	mainObj := tokenArr{}

	//блок работы
	//конвертируем в строки
	for scanner.Scan() {
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return
		}
		mainObj.main = append(mainObj.main, lines)
	}

	//выделяем сортируемую и несортьируемую часть по колонке
	mainObj.getFreeSort(*colNum - 1)

	//выюираем метод сортировки
	if *revSort {
		if *numSort {
			sort.Sort(ByNumberAsc{&mainObj})
		} else {
			sort.Sort(ByLetterAsc{&mainObj})
		}
	} else {
		if *numSort {
			sort.Sort(ByNumberDesc{&mainObj})
		} else {
			sort.Sort(ByLetterDesc{&mainObj})
		}
	}

	//удаляем дуюликаты на уже отсоритрованном массиве потому что так проще
	if *dupSort {
		Unique(&mainObj)
	}

	for _, v := range mainObj.main {
		fmt.Println(string(v))
	}

}
