package internal

import (
	"fmt"
	"regexp"
	"strings"
)

// Match - Новый тип. функция поиска совпадения
type Match func(line string) bool

// BildMatcher - создаёт функцию поиска совпадения
// flags - структра с флагами, pattern который необходимо найти
func BildMatcher(flags Flags, pattern string) (Match, error) {
	realPattern := pattern
	// fixed pattern
	if flags.Fixed { // -F
		if flags.IgnoreCase { // -i
			realPattern = strings.ToLower(realPattern)
		}
		return fixedFlagFunc(realPattern, flags), nil
	}

	if flags.IgnoreCase { // -i
		realPattern = "(?i)" + realPattern
	}
	// regexp
	reg, err := regexp.Compile(realPattern)
	if err != nil {
		return nil, fmt.Errorf("Ошибка компиляции регулярного выражения: %q, %w", pattern, err)
	}
	return func(line string) bool {
		ok := reg.MatchString(line)
		if flags.Invert {
			return !ok
		}
		return ok
	}, nil
}

func fixedFlagFunc(pattern string, flags Flags) Match {
	return func(line string) bool {
		str := line
		if flags.IgnoreCase {
			str = strings.ToLower(str)
		}
		ok := strings.Contains(str, pattern)
		if flags.Invert {
			return !ok
		}
		return ok
	}
}
