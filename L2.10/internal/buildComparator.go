package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type LessFunc func(a, b string) bool

func BuildLess(flags LineArgs) (LessFunc, error) {
	less := LexLess // базово — лексикографически

	// из nMh не может быть только 1
	if flags.N {
		less = NumericPrefLess
	}
	if flags.H {
		less = HumanLess
	}
	if flags.B {
		// fmt.Println("ЗАШЛО В ФЛАГ -b")
		less = WrapIgnoreTrailingBlanks(less)
	}
	if flags.K > 0 {
		less = WrapKeyColumn(less, flags.K)
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

// NumericPrefLess сравнивает строки как числа
// Если числа равны, то сравниваем исходные строки лексикографически.
func NumericPrefLess(a, b string) bool {
	aTrim := strings.TrimLeft(a, " ") // в GNU sort не смотрятся пробелы если числа
	bTrim := strings.TrimLeft(b, " ") //

	av, _ := leadingNumberOrZero(aTrim)
	bv, _ := leadingNumberOrZero(bTrim)

	if av < bv {
		return true
	}
	if av > bv {
		return false
	}

	// tie-breaker: лексикографически исходные строки
	return a < b
}

// возвращает число и индекс последней цифры (-1, если нет цифр)
// исправь чтобы обрабатывало точку в начале (.123, -.123, +.123 )
func leadingNumberOrZero(s string) (float64, int) {
	i := 0
	n := len(s)
	if n == 0 {
		return 0, -1
	}

	// +-
	if s[i] == '+' || s[i] == '-' {
		i++
		if i >= n {
			return 0, -1
		}
	}

	startDigits := i

	// целая часть или 0
	for i < n && isDigit(s[i]) {
		i++
	}

	hasIntDigits := i > startDigits

	// после точки, если она есть
	if i < n && s[i] == '.' {
		j := i + 1
		// числа после точки
		for j < n && isDigit(s[j]) {
			j++
		}
		hasFracDigits := j > i+1
		if hasIntDigits || hasFracDigits {
			i = j
			hasIntDigits = hasIntDigits || hasFracDigits // .123 - валидно (0.123)
		}
	}
	// нет цифр в начале — нет числа
	if !hasIntDigits {
		return 0, -1
	}
	// переводим найденное число в строку
	numStr := s[:i]
	v, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, 0
	}
	return v, i
}
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func MonthLess(less LessFunc) LessFunc {
	return func(a, b string) bool {
		// замена первого вхождения названия месяца на его номер
		aMonth := parseMonth(a)
		bMonth := parseMonth(b)
		return less(aMonth, bMonth)
	}
}

var month = map[string]string{"Jan": "1", "jan": "1", "Feb": "2", "feb": "2", "Mar": "3", "mar": "3", "Apr": "4", "apr": "4",
	"May": "5", "may": "5", "Jun": "6", "jun": "6", "Jul": "7", "jul": "7", "Aug": "8", "aug": "8", "Sep": "9", "sep": "9",
	"Oct": "10", "oct": "10", "Nov": "11", "nov": "11", "Dec": "12", "dec": "12"}

func parseMonth(str string) string {
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

func HumanLess(a, b string) bool {
	aVal, indA := leadingNumberOrZero(a)
	bVal, indB := leadingNumberOrZero(b)

	if indA == -1 {
		aVal = findSuffixAndMultyply(aVal, a, indA)
	}
	if indB == -1 {
		bVal = findSuffixAndMultyply(bVal, b, indB)
	}
	if aVal < bVal {
		return true
	}
	if aVal > bVal {
		return false
	}

	// tie-breaker: лексикографически исходные строки
	return a < b
}

var suffix = map[byte]float64{'B': 1, 'b': 1, 'K': 1 << 10, 'k': 1 << 10, 'M': 1 << 20, 'm': 1 << 20, 'G': 1 << 30, 'g': 1 << 30}

// если есть суффикс B, K, M, G, то увеличь в нужное число раз число, которое перед суффиксом
func findSuffixAndMultyply(num float64, str string, indexLastDigit int) float64 {
	// нет символов потом
	if indexLastDigit == len(str)-1 {
		return num
	}

	if val, ok := suffix[str[indexLastDigit+1]]; ok {
		num *= val // умножаем, если есть суффикс
	}

	return num
}

func WrapIgnoreTrailingBlanks(less LessFunc) LessFunc {
	return func(a, b string) bool {
		aTrim := strings.TrimRight(a, " \t")
		// fmt.Printf("ЗАТРИМИЛАСЬ ПЕРВАЯ СТРОКА\nбыло:<%s>\nстало:<%s>\n", a, aTrim)
		bTrim := strings.TrimRight(b, " \t")
		// fmt.Printf("ЗАТРИМИЛАСЬ ВТОРАЯ СТРОКА\nбыло:<%s>\nстало:<%s>\n", a, aTrim)

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
