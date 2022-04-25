package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	reader := bufio.NewReader(os.Stdin)
	fields := flag.Int("f", 0, "выбрать поля (колонки)")
	delimiter := flag.String("d", "	", "использовать другой разделитель (базовый TAB)")
	separated := flag.Bool("s", false, "только строки с разделителем")
	//reader := bufio.NewReader(os.Stdin)
	flag.Parse()
	fmt.Println("delimiter: ", *delimiter)
	for {

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		tokens := strings.Split(input, *delimiter)
		wasprint := false
		for i, e := range tokens {
			if i >= *fields-1 {
				fmt.Print(e)
				wasprint = true
			}
		}
		if !wasprint && len(tokens) == 1 && !*separated {
			fmt.Print(input)
		}
	}
}
