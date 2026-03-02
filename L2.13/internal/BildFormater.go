package internal

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Formater - функция форматирования строки.
// Договоримся: пустая строка "" => "не выводить"
type Formater func(line string) (string, bool)

func BildFormater(flags Flags) (Formater, error) {
	formater := defaultFormatter

	if flags.Fields != "" {
		ff, err := fieldsFormater(flags)
		if err != nil {
			return nil, err
		}
		formater = ff
	}

	if flags.Separated {
		sep := flags.Delimiter
		if sep == "" {
			sep = "\t"
		}
		formater = separatedFormater(formater, sep)
	}

	return formater, nil
}

func defaultFormatter(line string) (string, bool) { return line, true }

func separatedFormater(form Formater, sep string) Formater {
	return func(line string) (string, bool) {
		if strings.Contains(line, sep) {
			return form(line)
		}
		return "", false
	}
}

func fieldsFormater(flags Flags) (Formater, error) {

	sep := flags.Delimiter
	if sep == "" {
		sep = "\t"
	}

	idxs, err := parseFields(flags.Fields)
	if err != nil {
		return nil, err
	}

	return func(line string) (string, bool) {
		if !strings.Contains(line, sep) {
			return line, true
		}
		parts := strings.Split(line, sep)

		// Собираем выбранные поля
		out := make([]string, 0, len(idxs))
		for _, idx := range idxs {
			if idx >= 0 && idx < len(parts) {
				out = append(out, parts[idx])
			}
		}

		if len(out) == 0 {
			return "", true
		}

		return strings.Join(out, sep), true
	}, nil
}

// parseFields парсит "1,3-5,7" -> []int{0,2,3,4,6}
func parseFields(spec string) ([]int, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("fields: пустое значение")
	}

	set := make(map[int]struct{}, 16)

	tokens := strings.Split(spec, ",")
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			return nil, fmt.Errorf("fields: пустой элемент в списке %s", spec)
		}

		if strings.Contains(tok, "-") {
			a, b, ok := strings.Cut(tok, "-")
			if !ok {
				return nil, fmt.Errorf("fields: неверный диапазон %q", tok)
			}
			a = strings.TrimSpace(a)
			b = strings.TrimSpace(b)

			start, err := strconv.Atoi(a)
			if err != nil {
				return nil, fmt.Errorf("fields: неверное число %q", a)
			}
			end, err := strconv.Atoi(b)
			if err != nil {
				return nil, fmt.Errorf("fields: неверное число %q", b)
			}
			if start <= 0 || end <= 0 {
				return nil, fmt.Errorf("fields: номера полей должны быть > 0: %q", tok)
			}
			if start > end {
				return nil, fmt.Errorf("fields: start > end в диапазоне %q", tok)
			}

			for i := start; i <= end; i++ {
				set[i-1] = struct{}{}
			}
			continue
		}

		// одиночное число
		n, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("fields: неверное число %q", tok)
		}
		if n <= 0 {
			return nil, fmt.Errorf("fields: номера полей должны быть > 0: %q", tok)
		}
		set[n-1] = struct{}{}
	}

	out := make([]int, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Ints(out)
	return out, nil
}
