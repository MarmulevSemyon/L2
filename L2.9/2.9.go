package main

import (
	"fmt"
	"os"
	"unicode"
)

// Написать функцию Go, осуществляющую примитивную распаковку строки, содержащей повторяющиеся символы/руны.

// Примеры работы функции:
// Вход: "a4bc2d5e"
// Выход: "aaaabccddddde"
// Вход: "abcd"
// Выход: "abcd" (нет цифр — ничего не меняется)
// Вход: "45"
// Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)
// Вход: ""
// Выход: "" (пустая строка -> пустая строка)

// Дополнительное задание
// Поддерживать escape-последовательности вида \:
// Вход: "qwe\4\5"
// Выход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)
// Вход: "qwe\45"
// Выход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)
// Требования к реализации
// Функция должна корректно обрабатывать ошибочные случаи (возвращать ошибку, например, через error), и проходить unit-тесты.
// Код должен быть статически анализируем (vet, golint).
func main() {
	packedString := "r"

	err, unpackedString := unpackString(packedString)
	if err != 0 {
		os.Exit(1)
	}
	fmt.Println(unpackedString)

}

func unpackString(packedString string) (int, string) {
	arrRune := []rune(packedString)

	escapingFlag := false
	prev := arrRune[0]
	resRune := make([]rune, len(packedString))
	count := 0

	if len(arrRune) == 0 {
		return 0, ""
	}
	if unicode.IsDigit(arrRune[0]) {
		return 1, ""
	}

	for i := 0; i < len(arrRune); i++ {
		if i != 0 {
			prev = arrRune[i-1]
		}
		count = 1
		if prev == '\\' {
			escapingFlag = true
			if i == len(arrRune)-1 {
			} else {
				resRune = addingToRes(resRune, prev, count, &escapingFlag)
			}
			resRune = addingToRes(resRune, arrRune[i], count, &escapingFlag)
			break
		}

		if unicode.IsDigit(arrRune[i]) {
			if unicode.IsDigit(prev) && !escapingFlag {
				return 1, ""
			}
			if prev == '\\' {
				escapingFlag = true
				continue
			}
			count = int(arrRune[i] - '0')
		}

		resRune = addingToRes(resRune, prev, count, &escapingFlag)
	}
	return 0, string(resRune)
}

func addingToRes(res []rune, char rune, count int, escapingFlag *bool) []rune {

	if unicode.IsDigit(char) && !(*escapingFlag) {
		return res
	}

	addedArr := make([]rune, count)
	for i := 0; i < count; i++ {
		addedArr[i] = char
	}
	res = append(res, addedArr...)
	*escapingFlag = false
	return res
}
