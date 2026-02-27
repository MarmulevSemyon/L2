package internal

import (
	"bufio"
	"fmt"
	"io"
)

type ctxLine struct {
	number int
	text   string
}

// PrintGrep печатает нужные строки с учетом флагов
func PrintGrep(r io.Reader, match Match, flags Flags) error {
	sc := bufio.NewScanner(r)

	before := make([]ctxLine, 0, flags.Before)
	afterLeft := 0

	lastPrinted := 0
	lineNumber := 0

	for sc.Scan() {
		lineNumber++
		line := sc.Text()

		isMatch := match(line)
		if isMatch {
			flushBefore(before, &lastPrinted, flags.LineNum)

			// печатаем текущую строку (если ещё не печатали)
			if lineNumber > lastPrinted {
				printLine(flags.LineNum, lineNumber, line)
				lastPrinted = lineNumber
			}

			afterLeft = flags.After
			continue
		}

		// если нет совпадения, то печатаем контекст after
		if afterLeft > 0 {
			if lineNumber > lastPrinted {
				printLine(flags.LineNum, lineNumber, line)
				lastPrinted = lineNumber
			}
			afterLeft--
			continue
		}

		// иначе копим элик (контекст before)
		if flags.Before > 0 {
			if len(before) == flags.Before {
				copy(before, before[1:])
				before = before[:flags.Before-1]
			}
			before = append(before, ctxLine{number: lineNumber, text: line})
		}
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func flushBefore(before []ctxLine, lastPrinted *int, flagN bool) {
	for _, line := range before {
		// не печатать повторно
		if line.number <= *lastPrinted {
			continue
		}
		printLine(flagN, line.number, line.text)
		*lastPrinted = line.number
	}
	before = before[:0]
}

func printLine(flagN bool, num int, line string) {
	if flagN {
		fmt.Printf("%d:%s\n", num, line)
	} else {
		fmt.Println(line)
	}
}

// CountOfMatch считает количество строк, содержащих паттерн, в файле
// возвращает количество строк и ошибку
func CountOfMatch(match Match, r io.Reader) (int, error) {
	count := 0
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if match(line) {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
