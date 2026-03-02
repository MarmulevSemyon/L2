package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFieldsSinglRangeAndDupl(t *testing.T) {
	t.Parallel()

	idxs, err := parseFields("1,3-5,3, 5 ")
	require.NoError(t, err)
	// 1 -> 0, 3-5 -> 2,3,4, duplicate 3/5 ignored
	assert.Equal(t, []int{0, 2, 3, 4}, idxs)
}

func TestParseFieldsTrimsSpaces(t *testing.T) {
	t.Parallel()

	idxs, err := parseFields("  2 , 4 - 6 ")
	require.NoError(t, err)
	assert.Equal(t, []int{1, 3, 4, 5}, idxs)
}

func TestParseFieldsEmptySpec(t *testing.T) {
	t.Parallel()

	_, err := parseFields("   ")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: пустое значение")
}

func TestParseFieldsEmptyToken(t *testing.T) {
	t.Parallel()

	_, err := parseFields("1,,2")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: пустой элемент")
}

func TestParseFieldsNotNumber(t *testing.T) {
	t.Parallel()

	_, err := parseFields("1,a")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: неверное число")
}

func TestParseFieldsZeroOrNegative(t *testing.T) {
	t.Parallel()

	_, err := parseFields("0,1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "должны быть > 0")

	_, err = parseFields("-2") // будет считаться как диапазон
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: неверное число")
}

func TestParseFieldsRangeEndStart(t *testing.T) {
	t.Parallel()

	_, err := parseFields("5-3")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "start > end")
}

func TestParseFieldRangeBadNumber(t *testing.T) {
	t.Parallel()

	_, err := parseFields("a-3")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: неверное число")

	_, err = parseFields("1-b")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields: неверное число")
}

func TestFieldsFormaterDefaultDelimiterTab(t *testing.T) {
	t.Parallel()

	f, err := fieldsFormater(Flags{Fields: "1,3"})
	require.NoError(t, err)

	got, ok := f("a\tb\tc\td")
	require.True(t, ok)
	assert.Equal(t, "a\tc", got)
}

func TestFieldsFormaterCustomDelimiter(t *testing.T) {
	t.Parallel()

	f, err := fieldsFormater(Flags{Fields: "2-4", Delimiter: ","})
	require.NoError(t, err)

	got, ok := f("a,b,c,d,e")
	require.True(t, ok)
	assert.Equal(t, "b,c,d", got)
}

func TestFieldsFormaterOutOfRangeIgnored(t *testing.T) {
	t.Parallel()

	f, err := fieldsFormater(Flags{Fields: "10,2", Delimiter: ","})
	require.NoError(t, err)

	got, ok := f("a,b,c")
	require.True(t, ok)
	assert.Equal(t, "b", got)
}

func TestFieldsFormaterNoSelectedFieldsReturnsEmpty(t *testing.T) {
	t.Parallel()

	f, err := fieldsFormater(Flags{Fields: "10", Delimiter: ","})
	require.NoError(t, err)

	got, ok := f("a,b,c")
	require.True(t, ok)
	assert.Equal(t, "", got)
}

func TestSeparatedFormaterFiltersLinesWithoutDelimiter(t *testing.T) {
	t.Parallel()

	inner := func(s string) (string, bool) { return "X:" + s, true }
	f := separatedFormater(inner, ",")

	got, ok := f("a,b")
	require.True(t, ok)
	assert.Equal(t, "X:a,b", got)

	got, ok = f("ab")
	require.False(t, ok)
	assert.Equal(t, "", got)
}

func TestBildFormaterDefNoFlags(t *testing.T) {
	t.Parallel()

	f, err := BildFormater(Flags{})
	require.NoError(t, err)

	got, ok := f("hello")
	require.True(t, ok)
	assert.Equal(t, "hello", got)
}

func TestBildFormater_SeparatedAndFields(t *testing.T) {
	t.Parallel()

	f, err := BildFormater(Flags{Fields: "2", Delimiter: ",", Separated: true})
	require.NoError(t, err)

	got, ok := f("a,b,c")
	require.True(t, ok)
	assert.Equal(t, "b", got)

	got, ok = f("abc")
	require.False(t, ok)
	assert.Equal(t, "", got)
}

func TestBildFormaterBadFieldsSpec(t *testing.T) {
	t.Parallel()

	_, err := BildFormater(Flags{Fields: "1,,2"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fields:")
}
