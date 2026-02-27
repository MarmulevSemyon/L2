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

func TestParseLine(t *testing.T) {
	flagsErr := [][]string{{"grep", "-i", "-v", "-F", "-n", "-A", "1", "-B", "-1", "C", "12"},
		{"grep", "-ivFn", "-A", "1", "-B", "-1"},
		{"grep", "-ivFn", "-B", "-1", "C", "12"},
		{"grep", "-ivF", "-A", "1", "-B", "-1", "C", "12"},
		{"grep", "-ivn", "-B", "-1"},
		{"grep", "-iv", "-B", "-1"},
		{"grep", "-i", "-B", "-1"},
	}

	empty := Flags{}
	for i := range flagsErr {
		actual, _, err := ParseLine(flagsErr[i])
		assert.NotNil(t, err)
		assert.Equal(t, empty, actual)
	}

	flags1 := [][]string{{"grep", "-i", "-v", "-F", "-n", "-A", "1", "-c", "-B", "1", "-C", "12"},
		{"grep", "-ivFnc", "-A", "1", "-B", "1", "-C", "12"},
		{"grep", "-ivncF", "-A", "1", "-C", "12", "-B", "1"},
		{"grep", "-vicFn", "-C", "12", "-B", "1", "-A", "1"},
		{"grep", "-vcFin", "-C", "12", "-A", "1", "-B", "1"},
		{"grep", "-cvFin", "--context", "12", "--after-context", "1", "--before-context", "1"},
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

// FlagsM := [][]string{{"grep", "-c", "-r", "-M", "-u", "-b", "-k", "12"},
// 	{"grep", "-crMubk12"},
// 	{"grep", "-crMubk=12"},
// 	{"grep", "-crMubk", "12"},
// 	{"grep", "-crMub", "--key", "12"},
// 	{"grep", "--key", "12", "-crMub"},
// }
// expectedM := Flags{
// 	K: 12,
// 	N: false,
// 	R: true,
// 	U: true,
// 	M: true,
// 	B: true,
// 	C: true,
// 	H: false,
// }
// for _, v := range FlagsM {
// 	actual, _, err := ParseLine(v)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedM, actual)
// }

// FlagsH := [][]string{{"grep", "-c", "-r", "-h", "-u", "-b", "-k", "12"},
// 	{"grep", "-crhubk12"},
// 	{"grep", "-crhubk=12"},
// 	{"grep", "-crhubk", "12"},
// 	{"grep", "-crhub", "--key", "12"},
// 	{"grep", "--key", "12", "-crhub"},
// }
// expectedH := Flags{
// 	K: 12,
// 	N: false,
// 	R: true,
// 	U: true,
// 	M: false,
// 	B: true,
// 	C: true,
// 	H: true,
// }
// for _, v := range FlagsH {
// 	actual, _, err := ParseLine(v)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedH, actual)
// }
