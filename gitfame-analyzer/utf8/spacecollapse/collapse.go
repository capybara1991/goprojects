//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))

	inSpace := false

	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && w == 1 {
			r = utf8.RuneError
		}
		i += w

		if unicode.IsSpace(r) {
			inSpace = true
			continue
		}

		if inSpace {
			b.WriteByte(' ')
			inSpace = false
		}
		b.WriteRune(r)
	}

	if inSpace {
		b.WriteByte(' ')
	}
	return b.String()
}
