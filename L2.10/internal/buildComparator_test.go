package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericPrefLess(t *testing.T) {
	strs := []string{"-098",
		" апвуп",
		"0000",
		"0ю123",
		"qweqfsd",
		".123",
		"0.123",
		"11",
		" 98.",
		"98",
		"103",
		" 123",
		"123qaqeq",
		"1024A",
		"1024K",
		"100000,213215"}

	for i := 0; i < len(strs)-1; i++ {
		actual := numericPrefLess(strs[i], strs[i+1])
		// fmt.Printf("strs[i] =\t<%s>\nstrs[i+1] =\t<%s>\nbool =\t%v\n", strs[i], strs[i+1], actual)
		assert.Equal(t, true, actual)
	}
}

func TestParseHuman(t *testing.T) {
	strs := []string{"-2M",
		"-1K",
		"0",
		"abc",
		"строка",
		"я не понимаю как этот флаг сортирует 2000 и 1K",
		"я сделал как это работает в wsl and mint",
		"1",
		"12",
		"100B",
		"999",
		"1025",
		"0.5K",
		"1K",
		"1k",
		"2K",
		"10K",
		"512K",
		"1024K",
		"0.25M",
		"1.0M",
		"1M",
		"1.5M",
		"2M",
		"10M",
		"512M",
		"1030M",
		"1.0G",
		"1G",
		"1.5G",
		"2G",
		"10G",
		"512G",
		"1024G",
	}

	for i := 0; i < len(strs)-1; i++ {
		actual := humanLess(strs[i], strs[i+1])
		// fmt.Printf("strs[i] =\t%s\nstrs[i+1] =\t%s\nbool =\t%v\n", strs[i], strs[i+1], actual)
		assert.Equal(t, true, actual)
	}
}
func TestParseMonth(t *testing.T) {
	strs := []string{"sort", "-.123c", "123", "-qwe", ".qjanwe-M", "jan-.qwe-b", "feb-0.-h", "mar-", ".Mar", "May", "", "jun", "-"}

	expect := []string{"0sort", "0-.123c", "0123", "0-qwe", "0.qjanwe-M", "1-.qwe-b", "2-0.-h", "3-", "0.Mar", "5", "0", "6", "0-"}
	for i := range strs {
		actual := parseMonth(strs[i])
		assert.Equal(t, expect[i], actual)
	}
}
func TestTrim(t *testing.T) {
	strs := []string{"sort   ", "-.123c", ".-123r   \t\t", "123b\t", "   123B", "\t-Gqwe", ".Mqwe-M\t   ", "-.q   we"}

	expect := []string{"sort", "-.123c", ".-123r", "123b", "   123B", "\t-Gqwe", ".Mqwe-M", "-.q   we"}
	for i := range strs {
		actual := strings.TrimRight(strs[i], " \t")
		assert.Equal(t, expect[i], actual)
	}
}

func TestTGetValueByKIndex(t *testing.T) {
	strs := []string{"sort\tqwe", "-.123c\t", ".-123r\t \t123", "123b\t.\t", "   123B", "\t-Gqwe", ".Mqwe-M\t   ", "-.q   we\t123\t"}

	expect := []string{"qwe", "", " ", ".", "", "-Gqwe", "   ", "123"}
	for i := range strs {
		actual := getValueByKIndex(strs[i], 2)
		assert.Equal(t, expect[i], actual)
	}
}
