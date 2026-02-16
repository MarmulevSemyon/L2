package internal

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type LessFunc func(a, b string) bool

func BuildLess(flags LineArgs) (LessFunc, error) {
	less := LexLess // базово — лексикографически

	if flags.B {
		less = WrapIgnoreTrailingBlanks(less)
	}
	if flags.K > 0 {
		less = WrapKeyColumn(less, flags.K)
	}
	// одновремменно nMh не может быть
	if flags.N {
		less = NumericLess(less)
	}
	if flags.H {
		less = HumanLess(less)
	}
	if flags.M {
		less = MonthLess(less)
	}
	if flags.R {
		less = WrapReverse(less)
	}
	return less, nil
}

func LexLess(a, b string) bool {
	return a < b
}

// Лексикографическое сравнение, но строка "безчисел" считаются как "0безчисел"
func NumericLess(less LessFunc) LessFunc {
	return func(a, b string) bool {
		aNum := extendZeroIfString(a)
		bNum := extendZeroIfString(b)
		return less(aNum, bNum)
	}
}

// Ставим ноль в начало, если не число
func extendZeroIfString(str string) string {
	strRune := []rune(str)
	if len(str) == 0 {
		strRune = []rune{'0'}
	} else if len(str) == 1 && unicode.IsGraphic(strRune[0]) {
		strRune = append([]rune{'0'}, strRune...)

	} else if !unicode.IsDigit(strRune[0]) && !(strRune[0] == '-' && unicode.IsDigit(strRune[1])) {
		strRune = append([]rune{'0'}, strRune...)
	}

	return string(strRune)
}

func MonthLess(less LessFunc) LessFunc {
	return func(a, b string) bool {
		// замена первого вхождения названия месяца на его номер
		aMonth := parseMonth(a)
		bMonth := parseMonth(b)
		return less(aMonth, bMonth)
	}
}
func parseMonth(str string) string {
	month := map[string]string{"Jan": "1", "jan": "1", "Feb": "2", "feb": "2", "Mar": "3", "mar": "3", "Apr": "4", "apr": "4",
		"May": "5", "may": "5", "Jun": "6", "jun": "6", "Jul": "7", "jul": "7", "Aug": "8", "aug": "8", "Sep": "9", "sep": "9",
		"Oct": "10", "oct": "10", "Nov": "11", "nov": "11", "Dec": "12", "dec": "12"}

	res := []rune(str)

	if len(res) < 3 {
		res = append([]rune("0"), res...)
	} else if val, ok := month[string(res[:3])]; ok {
		res = append([]rune(val), res[3:]...)
	} else {
		res = append([]rune("0"), res...)
	}

	return string(res)
}

func HumanLess(less LessFunc) LessFunc {
	return func(a, b string) bool {
		aHuman := parseHuman(a)
		bHuman := parseHuman(b)
		return less(aHuman, bHuman)
	}
}

// если есть суффикс B, K, M, G, то увеличь в нужное число раз число, которое перед суффиксом
func parseHuman(str string) string {
	suffix := map[rune]int{'B': 1, 'b': 1, 'K': 1 << 10, 'k': 1 << 10, 'M': 1 << 20, 'm': 1 << 20, 'G': 1 << 30, 'g': 1 << 30}
	runeStr := []rune(str)
	startChar := len(runeStr)

	for i, v := range runeStr {
		if unicode.IsDigit(v) || (v == '-' && i == 0 && len(str) > 1 && unicode.IsDigit(runeStr[i+1])) {
			continue
		} else if v == '.' && i != 0 && unicode.IsDigit(runeStr[i-1]) {
			continue
		} else {
			startChar = i
			break
		}
	}

	if startChar == 0 {
		res := []rune(str)
		res = append([]rune("0"), res...)
		return string(res)
	}

	chislo, _ := strconv.ParseFloat(string(runeStr[:startChar]), 64)
	if startChar == len(runeStr) { // чисто число
		return str
	}

	if val, ok := suffix[runeStr[startChar]]; ok {
		chislo *= float64(val) // умножаем, если есть суффикс
	} else {
		return str
	}

	strChislo := strconv.FormatFloat(chislo, 'f', -1, 64) // переводим в строку
	if startChar == len(runeStr)-1 {                      // если суффикс последний эл-т в строке
		return strChislo
	}

	res := append([]rune(strChislo), runeStr[startChar+1:]...)
	fmt.Printf("перевели уже,\nстрока: %s\nначало символов %d\nчто будет после числа: %s\n ", string(res), startChar, string(runeStr[startChar:]))

	return string(res)
}

func WrapIgnoreTrailingBlanks(less LessFunc) LessFunc {
	return func(a, b string) bool {
		aTrim := strings.Trim(a, " \t")
		bTrim := strings.Trim(b, " \t")
		return less(aTrim, bTrim)
	}
}

// TODO: если равны то сравнивать по всей строке
func WrapKeyColumn(less LessFunc, key int) LessFunc {
	return func(a, b string) bool {
		aK := getValueByKIndex(a, key)
		bK := getValueByKIndex(b, key)
		return less(aK, bK)
	}
}
func getValueByKIndex(str string, ind int) string {
	res := strings.Split(str, "\t")
	fmt.Printf("")
	if len(res) < ind-1 {
		return ""
	}
	return res[ind-1]
}

func WrapReverse(less LessFunc) LessFunc {
	return func(a, b string) bool {
		return less(b, a)
	}
}
