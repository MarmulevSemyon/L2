package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// type Flags struct {
// After  int 		// -A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).
// Before int 		// -B N — вывести N строк до каждой найденной строки.
// Context int 	// -C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).
// Count bool 		// -c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).
// IgnoreCase bool // -i — игнорировать регистр.
// Invert bool 	// -v — инвертировать фильтр: выводить строки, не содержащие шаблон.
// Fixed bool 		// -F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).
// LineNum bool 	// -n — выводить номер строки перед каждой найденной строкой.

// 	CPUProfile string
// 	MemProfile string
// 	Trace      string
// }

func TestParseLineErr1(t *testing.T) {
	flagsErr := [][]string{{"-i", "-v", "-F", "-n", "-A", "1", "-B", "-1", "-C", "12", "pattern", "file"},
		{"-ivFn", "-A", "1", "-B", "-1", "pattern", "file"},
		{"-ivFn", "-B", "-1", "-C", "12", "pattern", "file"},
		{"-ivF", "-A", "1", "-B", "-1", "-C", "12", "pattern", "file"},
		{"-ivn", "-B", "-1", "pattern", "file"},
		{"-iv", "-B", "-1", "pattern", "file"},
		{"-i", "-B", "-1", "pattern", "file"},
	}

	empty := Flags{}
	for i := range flagsErr {
		actual, _, err := ParseLine(flagsErr[i])
		assert.NotNil(t, err)
		assert.Equal(t, empty, actual)
	}
}
func TestParseLineErr2(t *testing.T) {
	flagsErr := [][]string{{"-i", "-v", "-F", "-n", "-A", "1", "-C", "12", "file"},
		{"-ivFn", "-A", "1", "pattern"},
		{"-ivFn", "-C", "12", "file"},
		{"-ivF", "-A", "1", "-C", "12", "pattern"},
		{"-ivn", "file"},
		{"-iv", "pattern"},
		{"-i", "file"},
	}
	empty := Flags{}
	for i := range flagsErr {
		actual, _, err := ParseLine(flagsErr[i])
		assert.NotNil(t, err)
		assert.Equal(t, empty, actual)

	}
}

func TestParseLine(t *testing.T) {
	flags1 := [][]string{{"-i", "-v", "-F", "-n", "-A", "1", "-c", "-B", "1", "-C", "12", "pattern", "file"},
		{"-ivFnc", "-A", "1", "-B", "1", "-C", "12", "pattern", "file"},
		{"-ivncF", "-A", "1", "-C", "12", "-B", "1", "pattern", "file"},
		{"-vicFn", "-C", "12", "-B", "1", "-A", "1", "pattern", "file"},
		{"-vcFin", "-C", "12", "-A", "1", "-B", "1", "pattern", "file"},
		{"-cvFin", "--context", "12", "--after-context", "1", "--before-context", "1", "pattern", "file"},
	}
	expected1 := Flags{
		After:      1,
		Before:     1,
		Context:    12,
		Count:      true,
		IgnoreCase: true,
		Invert:     true,
		Fixed:      true,
		LineNum:    true,
	}
	for i := range flags1 {
		actual, _, err := ParseLine(flags1[i])
		assert.Nil(t, err)
		assert.Equal(t, expected1, actual)
	}
}
func TestParseLineCtx(t *testing.T) {
	flags1 := [][]string{{"-B", "1", "-C", "12", "pattern", "file"},
		{"-A", "1", "-C", "12", "pattern", "file"},
		{"-A", "1", "pattern", "file"},
		{"-B", "1", "pattern", "file"},
		{"-C", "12", "pattern", "file"},
	}
	expected := []Flags{
		{
			After:   12,
			Before:  1,
			Context: 12,
		},
		{
			After:   1,
			Before:  12,
			Context: 12,
		},
		{
			After:   1,
			Before:  0,
			Context: 0,
		},
		{
			After:   0,
			Before:  1,
			Context: 0,
		},
		{
			After:   12,
			Before:  12,
			Context: 12,
		},
	}
	for i := range flags1 {
		actual, _, err := ParseLine(flags1[i])
		assert.Nil(t, err)
		assert.Equal(t, expected[i], actual)
	}
}
