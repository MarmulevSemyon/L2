package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PrintIsSorted возвращаяет сообщение об отсортированности файла file согласно функции сравнения
func PrintIsSorted(file *os.File, less LessFunc) (string, error) {

	br := bufio.NewReader(file)

	var prev string
	hasPrev := false
	countStr := 0

	for {
		current, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		current = strings.TrimRight(current, "\r\n")
		if !hasPrev {
			hasPrev = true
			prev = current
		}
		countStr++

		if !less(prev, current) && !equalByLess(prev, current, less) {
			fmt.Printf("prev:\t<%s>\ncurrent:\t<%s>", prev, current)
			fmt.Println(less(prev, current))

			res := fmt.Sprintf("Не отсортировано после строчки №%d:\n\t%s\n", countStr, current)
			return res, nil
		}

		prev = current
		if err == io.EOF {
			break
		}
	}

	return "Файл отсортирован!", nil
}
