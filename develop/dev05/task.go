package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func regFreeContains(h, n []byte) bool {
	return bytes.Contains(bytes.ToLower(h), bytes.ToLower(n))
}

//это простой grep

func main() {

	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "(количество строк)")
	ignore := flag.Bool("i", false, "(игнорировать регистр)")
	invert := flag.Bool("v", false, "вместо совпадения, исключать")
	fixed := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	line := flag.Bool("n", false, "печатать номер строки")
	flag.Parse()

	args := flag.Args()

	file, err := os.Open(args[1])
	defer file.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	obj := [][]byte{}
	needle := []byte(args[0])
	printThis := false
	lineNum := 0
	if *context > 0 {
		*before = *context
		*after = *context
	}

	//загоняем файл в структуру
	for scanner.Scan() {
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return
		}
		obj = append(obj, lines)
	}

	printable := make(map[int]struct{})
	if *ignore {
		bytes.ToLower(needle)
	}
	for i, v := range obj {
		if *ignore {
			v = bytes.ToLower(v)
		}
		if *fixed {
			printThis = bytes.Contains(v, needle)
		} else {
			printThis, _ = regexp.Match(string(needle), v)
		}
		if *invert {
			printThis = !printThis
		}
		if printThis {
			lineNum++
			if !*count {
				printable[i] = struct{}{}
				b, a := *before, *after

				for b > 0 {
					if i-b >= 0 {
						printable[i-b] = struct{}{}
					}
					b--
				}
				for a > 0 {
					if i+b < len(obj) {
						printable[i+a] = struct{}{}
					}
					a--
				}
			}
		}
	}
	if *count {
		fmt.Println(lineNum)
	} else {
		var keys []int
		for k := range printable {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			if *line {
				fmt.Print(k+1, ":")
			}
			fmt.Print(string(obj[k]))
			fmt.Print("\n")
		}
	}
}
