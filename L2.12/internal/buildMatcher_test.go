package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBildMatcherFixed1(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      true,
		IgnoreCase: false,
		Invert:     false,
	}, "abc")

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, m("xxabcxx"), "fixed должен искать подстроку")
	assert.False(t, m("xxabxx"), "нет подстроки abc")
	assert.False(t, m("ABC"), "без -i регистр важен")
}
func TestBildMatcherFixed2(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      true,
		IgnoreCase: false,
		Invert:     false,
	}, "a.c")

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, m("xxa.cyy"), "должно совпасть буквально 'a.c'")
	assert.False(t, m("xxabcyy"), "не должно трактовать '.' как любой символ")
}
func TestBildMatcherFixedIgnoreCase(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      true,
		IgnoreCase: true,
		Invert:     false,
	}, "AbC")

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, m("xxabcxx"))
	assert.True(t, m("xxABCxx"))
	assert.True(t, m("xxaBcxx"))
	assert.False(t, m("xxabxx"))
}

func TestBildMatcherFixedInvert(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      true,
		IgnoreCase: false,
		Invert:     true,
	}, "abc")

	if !assert.NoError(t, err) {
		return
	}

	assert.False(t, m("xxabcxx"), "-v инвертирует совпадение")
	assert.True(t, m("xxabxx"), "если не совпало — true при -v")
}

func TestBildMatcherRegexp(t *testing.T) {
	m, err := BildMatcher(Flags{}, "a.c")

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, m("abc"))
	assert.True(t, m("aXc"))
	assert.False(t, m("ac"))
}

func TestBildMatcherRegexpIgnoreCase(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      false,
		IgnoreCase: true,
		Invert:     false,
	}, "курс")

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, m("КУРС"), "должно игнорировать регистр")
	assert.True(t, m("курс"))
	assert.True(t, m("КуРс"))
	assert.False(t, m("другое"))
}

func TestBildMatcherRegexpInvert(t *testing.T) {
	m, err := BildMatcher(Flags{
		Fixed:      false,
		IgnoreCase: false,
		Invert:     true,
	}, "abc")

	if !assert.NoError(t, err) {
		return
	}

	assert.False(t, m("xxabcxx"))
	assert.True(t, m("zzz"))
}

func TestBildMatcherInvalidPattern(t *testing.T) {
	_, err := BildMatcher(Flags{
		Fixed:      false,
		IgnoreCase: false,
		Invert:     false,
	}, "(") // некорректная regexp

	assert.Error(t, err)
}
