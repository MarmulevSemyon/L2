package internal

import (
	"bufio"
	"fmt"
	"io"
)

// PrintCat печатает нужные строки с учетом флагов
func PrintCat(r io.Reader, formater Formater) error {
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		line := sc.Text()
		formatedLine, ok := formater(line)
		if !ok {
			continue
		}
		fmt.Println(formatedLine)
	}

	return sc.Err()
}
