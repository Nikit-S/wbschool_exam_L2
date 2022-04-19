package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
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

type outPut struct {
	before   [][]byte
	after    int
	needle   []byte
	lineNum  int
	totalNum int
	distance int
	buf      bytes.Buffer
}

/*
func (o *outPut) checkLine(line []byte) bool {
	for _, f := range o.checkers {
		if !f(line, o.needle) {
			return false
		}
	}
	return true
}*/

func regFreeContains(h, n []byte) bool {
	return bytes.Contains(bytes.ToLower(h), bytes.ToLower(n))
}

func main() {

	mainObj := &outPut{}
	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "(количество строк)")
	ignore := flag.Bool("i", false, "(игнорировать регистр)")
	invert := flag.Bool("v", false, "вместо совпадения, исключать")
	fixed := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	line := flag.Bool("n", false, "печатать номер строки")
	flag.Parse()

	fmt.Println("after: ", *after)
	fmt.Println("before: ", *before)
	fmt.Println("context: ", *context)
	fmt.Println("count: ", *count)
	fmt.Println("ignore: ", *ignore)
	fmt.Println("invert: ", *invert)
	fmt.Println("fixed: ", *fixed)
	fmt.Println("line: ", *line)

	args := flag.Args()

	fmt.Println("args: ", args)

	file, err := os.Open(args[1])
	defer file.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	mainObj.needle = []byte(args[0])
	if *context > 0 {
		*before = *context
		*after = *context
	}
	for scanner.Scan() {
		mainObj.lineNum++
		mainObj.distance++
		printThis := false
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return
		}
		linesDup := lines
		if !*ignore && !*fixed {
			printThis, err = regexp.Match(string(mainObj.needle), linesDup)
		}
		if *ignore {
			linesDup := bytes.ToLower(lines)
			mainObj.needle = bytes.ToLower(mainObj.needle)
			printThis, err = regexp.Match(string(mainObj.needle), linesDup)
		}
		if *fixed {
			printThis = bytes.Contains(linesDup, mainObj.needle)
		}
		if *invert {
			printThis = !printThis
		}
		if printThis {
			if *count {
				mainObj.totalNum++
			}

			if *before > 0 {
				fmt.Println(mainObj.before)
				for j := range mainObj.before {
					if mainObj.distance-1-mainObj.after-len(mainObj.before)+j < 0 {
						continue
					}
					fmt.Fprintf(&mainObj.buf, "%s\n", mainObj.before[j])
				}
				mainObj.before = [][]byte{}
			}

			if *line {
				fmt.Fprintf(&mainObj.buf, "%v:", mainObj.lineNum)
			}
			fmt.Fprintf(&mainObj.buf, "%s", lines)
			fmt.Fprintf(&mainObj.buf, "\n")
			mainObj.after = 0
			mainObj.distance = 0
		}

		if *after > 0 && !printThis {
			if mainObj.after < *after {
				fmt.Fprintf(&mainObj.buf, "%s\n", lines)
				mainObj.after++
			}
		}
		if *before > 0 && !printThis {
			if *before == 1 {
				mainObj.before = [][]byte{lines}
			} else if len(mainObj.before) < *before {
				mainObj.before = append(mainObj.before, lines)
			} else {
				mainObj.before = append(mainObj.before[1:], lines)
			}
		}
	}
	if *count {
		fmt.Printf("%v\n", mainObj.totalNum)
	} else {
		fmt.Print(mainObj.buf.String())
	}
}
