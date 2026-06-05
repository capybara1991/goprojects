//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	if input == "" {
		return ""
	}

	total := 0
	for i := 0; i < len(input); {
		r, w := utf8.DecodeRuneInString(input[i:])
		if r == utf8.RuneError && w == 1 {
			r = utf8.RuneError
		}
		total += utf8.RuneLen(r)
		i += w
	}

	var b strings.Builder
	b.Grow(total)

	for i := len(input); i > 0; {
		r, w := utf8.DecodeLastRuneInString(input[:i])
		if r == utf8.RuneError && w == 1 {
			r = utf8.RuneError
		}
		b.WriteRune(r)
		i -= w
	}

	return b.String()
}
