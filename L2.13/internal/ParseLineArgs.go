package internal

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Flags - структура котрая хранит флаги
type Flags struct {
	Fields    string //-f "fields" — указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.
	Delimiter string //-d "delimiter" — использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').
	Separated bool   //-s – (separated) только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).

	CPUProfile string
	MemProfile string
	Trace      string
}

// ParseLine парсит массив строк(командной строки), возвращает структуру Flags.
func ParseLine(args []string) (Flags, []string, error) {
	var flags Flags
	fs := pflag.NewFlagSet("grep", pflag.ContinueOnError)

	fs.StringVar(&flags.CPUProfile, "cpuprofile", "", "write CPU profile to file")
	fs.StringVar(&flags.MemProfile, "memprofile", "", "write memory profile to file")
	fs.StringVar(&flags.Trace, "trace", "", "write execution trace to file")

	fs.StringVarP(&(flags.Fields), "fields", "f", "", "указание номеров полей (колонок), которые нужно вывести.")
	fs.StringVarP(&(flags.Delimiter), "delimiter", "d", "\t", "использовать другой разделитель (символ). По умолчанию: '\\t'.")
	fs.BoolVarP(&(flags.Separated), "separated", "s", false, " только строки, содержащие разделитель.")

	if err := fs.Parse(args); err != nil {
		return Flags{}, nil, fmt.Errorf("Ошибка аргументов: %w", err)
	}

	remainingArgs := fs.Args()
	if err := validateArgs(remainingArgs); err != nil {
		return Flags{}, nil, fmt.Errorf("Ошибка аргументов: %w", err)
	}

	return flags, remainingArgs, nil
}

func validateArgs(arg []string) error {
	if len(arg) == 0 {
		return fmt.Errorf("не указан файл")
	}
	return nil
}
