package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
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

var key string
var numeric bool
var reverse bool
var unique bool

func init() {
	testing.Init()
	flag.StringVar(&key, "k", "0", "sort via a key column")
	flag.BoolVar(&numeric, "n", false, "compare according to string numerical value")
	flag.BoolVar(&reverse, "r", false, "reverse the result of comparisons")
	flag.BoolVar(&unique, "u", false, "output only the first of an equal run")
	flag.Parse()
}

// имплементация sort.Interface для лексикографической сортировки
type stringSort struct {
	sl  [][]string
	key int
}

func (s stringSort) Len() int      { return len(s.sl) }
func (s stringSort) Swap(i, j int) { s.sl[i], s.sl[j] = s.sl[j], s.sl[i] }
func (s stringSort) Less(i, j int) bool {
	iLen := len(s.sl[i])
	jLen := len(s.sl[j])
	iHasKey := iLen > s.key
	jHasKey := jLen > s.key

	// сравнение по элементу key и случай, когда в одном из слайсов элементов меньше, чем key
	switch {
	case iHasKey && jHasKey:
		if s.sl[i][s.key] < s.sl[j][s.key] {
			return true
		}
		if s.sl[i][s.key] > s.sl[j][s.key] {
			return false
		}
	case !iHasKey && jHasKey:
		return true
	case iHasKey && !jHasKey:
		return false
	}

	// элементы слайсов с индексом key равны либо key лежит вне длин обоих слайсов
	for k := 0; k < minInt(iLen, jLen); k++ {
		if s.sl[i][k] < s.sl[j][k] {
			return true
		}
		if s.sl[i][k] > s.sl[j][k] {
			return false
		}
	}

	// слайсы равны по элементам, индексы которых есть в обоих слайсах
	return iLen < jLen
}

// имплементация sort.Interface для сортировки чисел
type numericSort struct {
	sl  [][]string
	key int
}

func (s numericSort) Len() int      { return len(s.sl) }
func (s numericSort) Swap(i, j int) { s.sl[i], s.sl[j] = s.sl[j], s.sl[i] }
func (s numericSort) Less(i, j int) bool {
	iLen := len(s.sl[i])
	jLen := len(s.sl[j])
	iHasKey := iLen > s.key
	jHasKey := jLen > s.key

	// сравнение по элементу key и случай, когда в одном из слайсов элементов меньше, чем key
	switch {
	case iHasKey && jHasKey:
		if lessNum(s.sl[i][s.key], s.sl[j][s.key]) {
			return true
		}
		if lessNum(s.sl[j][s.key], s.sl[i][s.key]) {
			return false
		}
	case !iHasKey && jHasKey:
		return true
	case iHasKey && !jHasKey:
		return false
	}

	// элементы слайсов с индексом key равны либо key лежит вне длин обоих слайсов
	for k := 0; k < minInt(iLen, jLen); k++ {
		if lessNum(s.sl[i][k], s.sl[j][k]) {
			return true
		}
		if lessNum(s.sl[j][k], s.sl[i][k]) {
			return false
		}
	}

	// слайсы равны по элементам, индексы которых есть в обоих слайсах
	return iLen < jLen
}

// Функция lessNum сравнивает строки как числа и возвращает true,
// если i < j либо j является числом, а i - нет.
func lessNum(i, j string) bool {
	iFloat, iIsNum := stringToFloat(i)
	jFloat, jIsNum := stringToFloat(j)
	if jIsNum && (!iIsNum || iFloat < jFloat) {
		return true
	}
	return false
}

// Функция minInt возвращает меньшее из двух целых чисел
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Функция stringToFloat конвертирует строку в число и возвращает его,
// а также идикатор успешного выполнения.
// Причем stringToFloat("NaN") = 0, false.
func stringToFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) {
		return 0, false
	}
	return f, true
}

func readStrings(file string) (lines []string, err error) {
	var scanner *bufio.Scanner
	if file == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(f)
		defer f.Close()
	}

	if unique {
		uLines := make(map[string]struct{})
		for scanner.Scan() {
			line := scanner.Text()
			if _, ok := uLines[line]; !ok {
				uLines[line] = struct{}{}
				lines = append(lines, line)
			}
		}
	} else {
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	return lines, nil
}

func readStringSlices(file string, k int) (lines [][]string, err error) {
	var scanner *bufio.Scanner
	if file == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(f)
		defer f.Close()
	}

	if unique {
		uLines := make(map[string]struct{})
		for scanner.Scan() {
			line := strings.Fields(scanner.Text())
			var field string
			if k < len(line) {
				field = line[k]
			}
			if _, ok := uLines[field]; !ok {
				uLines[field] = struct{}{}
				lines = append(lines, line)
			}
		}
	} else {
		for scanner.Scan() {
			lineString := scanner.Text()
			lines = append(lines, strings.Fields(lineString))
		}
	}

	return lines, nil
}

func main() {

	keyProvided := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "k" {
			keyProvided = true
		}
	})

	file := flag.Arg(0)

	var k int
	var err error
	if keyProvided {
		k, err = strconv.Atoi(key)
		if err != nil || k < 1 {
			fmt.Fprintf(os.Stderr,
				"invalid number at field start: invalid count at start of '%s'\n", key)
			return
		}
		k--
	}

	if k == 0 && !numeric {
		lines, err := readStrings(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		var strSlice sort.StringSlice = lines
		if reverse {
			sort.Sort(sort.Reverse(strSlice))
		} else {
			sort.Strings(lines)
		}
		for _, line := range lines {
			fmt.Println(line)
		}
		return
	}

	lines, err := readStringSlices(file, k)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var data sort.Interface
	if numeric {
		data = numericSort{lines, k}
	} else {
		data = stringSort{lines, k}
	}
	if reverse {
		data = sort.Reverse(data)
	}
	sort.Sort(data)
	for _, line := range lines {
		fmt.Println(strings.Join(line, " "))
	}
}
