package main

import (
	"fmt"
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

type MyStruct struct{}

func main() {
	m := make(map[struct{}]int)
	m[struct{}{}] = 2
	m[struct{}{}] = 3
	fmt.Println()

}
