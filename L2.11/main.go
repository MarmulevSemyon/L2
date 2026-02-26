package main

import (
	"fmt"
	"slices"
	"strings"
)

func main() {
	sliceIn := []string{"пятёк", "тяПкё", "пяткё", "листОк", "стоЛик", "слиток", "стол"}
	fmt.Println(findAnagram(sliceIn))
}

// Итоговая сложность O(n*m + n*log(n)), где n - количество слов, m - среднее количество букв в слове
// второе слогаемое появляется из-за необходимости сортировать слайсы слов по возрастанию.
// Если все слова будут анаграммами друг друга, то нужно будет сортировать все n слов
func findAnagram(slice []string) map[string][]string {
	res := make(map[string][]string)
	resWithCounterInKey := make(map[[33]uint16][]string)
	for _, str := range slice { // O(n)
		var counterRunes [33]uint16
		str = strings.ToLower(str) // O(m)
		for _, r := range str {    // O(m)
			if r == 'ё' || r == 'Ё' {
				counterRunes[32]++
			} else {
				counterRunes[r-'а']++ // кириллическая а
			}
		}
		resWithCounterInKey[counterRunes] = append(resWithCounterInKey[counterRunes], str) // O(1) в среднем
	}
	// O(2n*m) = O(n*m)

	for _, resSlice := range resWithCounterInKey { // O(h), где h - кол-во срезов анаграмм, граничные случаи:
		// hmax = n/2, тк если в слайсе 1 эл-т, то его пропускают
		// hmin = 1, в этом случае все слова относятся к 1 слайсу анаграмм, тогда его просто надо отсортировать (nlogn)
		if len(resSlice) == 1 {
			continue
		}
		firstStr := resSlice[0]  // O(1)
		slices.Sort(resSlice)    // О(n/h * log(n/h)), n/h - количество элементов среза (в худшем случае n, но тогда h = 1)
		res[firstStr] = resSlice // O(1)
	}
	// Итоговая сложность O(n*m) + O(h)*О(n/h * log(n/h))

	// h = n/2: O(n*m + n/2*2*log(2)) = O(n*(m+1)) = O(n*m) или
	// h = 1:  O(n*m + n*log(n))
	return res
}
