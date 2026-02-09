package internal

import (
	"fmt"

	"github.com/spf13/pflag"
)

// LineArgs хранит значения GNU-флагов утилиты sort.
type LineArgs struct {
	K int
	N bool
	R bool
	U bool

	M bool // сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат.
	B bool // — игнорировать хвостовые пробелы (trailing blanks).
	C bool // — проверить, отсортированы ли данные; если нет, вывести сообщение об этом.
	H bool // — сортировать по числовому значению с учётом суффиксов (например, К = килобайт, М = мегабайт — человекочитаемые размеры).
}

// ParseLine парсит аргументы массив строк(командной строки) и возвращает структуру LineArgs.
func ParseLine(args []string) (LineArgs, error) {
	var flags LineArgs
	fs := pflag.NewFlagSet("sort", pflag.ContinueOnError)
	fs.IntVarP(&(flags.K), "key", "k", 0, "сортировать по столбцу (колонке) №N.")
	fs.BoolVarP(&(flags.N), "numeric-sort", "n", false, "сортировать по числовому значению (строки интерпретируются как числа).")
	fs.BoolVarP(&(flags.R), "reverse", "r", false, " сортировать в обратном порядке (reverse).")
	fs.BoolVarP(&(flags.U), "unique", "u", false, "не выводить повторяющиеся строки (только уникальные).")

	fs.BoolVarP(&(flags.M), "month-sort", "M", false, "ортировать по названию месяца (Jan, Feb, ... Dec).")
	fs.BoolVarP(&(flags.B), "ignore-leading-blanks", "b", false, "игнорировать хвостовые пробелы.")
	fs.BoolVarP(&(flags.C), "check", "c", false, "проверить, отсортированы ли данные; если нет, вывести сообщение об этом.")
	fs.BoolVarP(&(flags.H), "human-numeric-sort", "h", false, "сортировать по числовому значению с учётом суффиксов (например, К = килобайт, М = мегабайт — человекочитаемые размеры).")

	if err := fs.Parse(args); err != nil {
		return LineArgs{}, fmt.Errorf("Ошибка аргументов: %w", err)
	}
	if err := validate(flags); err != nil {
		return LineArgs{}, fmt.Errorf("Ошибка аргументов: %w", err)
	}

	return flags, nil
}

// Одновремменно не может быть двух из 'nMh'
func validate(flags LineArgs) error {
	setOfNMH := []byte{'-'}

	if flags.N {
		setOfNMH = append(setOfNMH, 'n')
	}
	if flags.M {
		setOfNMH = append(setOfNMH, 'M')
	}
	if flags.H {
		setOfNMH = append(setOfNMH, 'h')
	}
	if len(setOfNMH) != 1 {
		return fmt.Errorf("флаги '%s' не совместимы", setOfNMH)
	}
	return nil
}
