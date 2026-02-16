package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtendZeroIfString(t *testing.T) {
	strs := []string{"sort", "-.123c", ".-123r", "123", "-qwe", ".qwe-M", "-.qwe-b", "-0.-h", "-", ".", "", "0.123", "-0.123"}

	expect := []string{"0sort", "0-.123c", "0.-123r", "123", "0-qwe", "0.qwe-M", "0-.qwe-b", "-0.-h", "0-", "0.", "0", "0.123", "-0.123"}
	for i := range strs {
		actual := extendZeroIfString(strs[i])
		assert.Equal(t, expect[i], actual)
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

func TestParseHuman(t *testing.T) {
	strs := []string{"sort", "-.123c", ".-123r", "123b", "123B", "-Gqwe", ".Mqwe-M", "-.qwe", "-1.G-h", "-", ".", "", "0.123k", "-0.123M"}

	expect := []string{"0sort", "0-.123c", "0.-123r", "123", "123", "0-Gqwe", "0.Mqwe-M", "0-.qwe", "-1073741824-h", "0-", "0.", "0", "125.952", "-128974.848"}
	for i := range strs {
		actual := parseHuman(strs[i])
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

	expect := []string{"qwe", "", " ", ".", "", "", "   ", "123"}
	for i := range strs {
		actual := getValueByKIndex(strs[i], 2)
		assert.Equal(t, expect[i], actual)
	}
}
