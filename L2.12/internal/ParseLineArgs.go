package internal

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Flags struct {
    After  int 		// -A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).
    Before int 		// -B N — вывести N строк до каждой найденной строки.
    Context int 	// -C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).
    Count bool 		// -c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).
    IgnoreCase bool // -i — игнорировать регистр.
    Invert bool 	// -v — инвертировать фильтр: выводить строки, не содержащие шаблон.
    Fixed bool 		// -F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).
    LineNum bool 	// -n — выводить номер строки перед каждой найденной строкой.
	
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

	fs.IntVarP(&(flags.After), "after-context", "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё (контекст)")
	fs.IntVarP(&(flags.Before), "before-context", "B", 0, "вывести N строк до каждой найденной строки.")
	fs.IntVarP(&(flags.Context), "context", "C", 0, "вывести N строк контекста вокруг найденной строки (включает и до, и после).")

	fs.BoolVarP(&(flags.Count), "count", "c", false, "выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число)")
	fs.BoolVarP(&(flags.IgnoreCase), "ignore-case", "i", false, "игнорировать регистр")
	fs.BoolVarP(&(flags.Invert), "invert-match", "v", false, "выводить строки, не содержащие шаблон")
	fs.BoolVarP(&(flags.Fixed), "fixed-strings ", "F", false, "выполнять точное совпадение подстроки")
	fs.BoolVarP(&(flags.LineNum), "line-number", "n", false, "выводить номер строки перед каждой найденной строкой.")

	if err := fs.Parse(args); err != nil {
		return Flags{}, nil, fmt.Errorf("Ошибка аргументов: %w", err)
	}
	if err := validate(flags); err != nil {
		return Flags{}, nil, fmt.Errorf("Ошибка аргументов: %w", err)
	}
	remainingArgs := fs.Args()
	return flags, remainingArgs, nil
}

func validate(flags Flags) error{
	if flags.After < 0{
		return fmt.Errorf("%d не может быть количеством строчек", flags.After)
	}
	if flags.Before < 0{
		return fmt.Errorf("%d не может быть количеством строчек", flags.Before)
	}
	if flags.Context < 0{
		return fmt.Errorf("%d не может быть количеством строчек", flags.Context)
	}
	return nil
}