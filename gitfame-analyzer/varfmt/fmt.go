//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	if format == "" {
		return ""
	}

	vals := make([]string, len(args))
	done := make([]bool, len(args))

	get := func(i int) string {
		if i < 0 || i >= len(args) {
			return ""
		}
		if done[i] {
			return vals[i]
		}
		switch v := args[i].(type) {
		case string:
			vals[i] = v
		case []byte:
			vals[i] = string(v)
		case int:
			vals[i] = strconv.Itoa(v)
		case int64:
			vals[i] = strconv.FormatInt(v, 10)
		case uint:
			vals[i] = strconv.FormatUint(uint64(v), 10)
		case uint64:
			vals[i] = strconv.FormatUint(v, 10)
		case bool:
			if v {
				vals[i] = "true"
			} else {
				vals[i] = "false"
			}
		default:
			vals[i] = fmt.Sprint(v)
		}
		done[i] = true
		return vals[i]
	}

	var b strings.Builder
	b.Grow(len(format) + 8*len(args))

	pos := 0
	for i := 0; i < len(format); {
		if format[i] != '{' {
			b.WriteByte(format[i])
			i++
			returnToContinue := false
			_ = returnToContinue
			continue
		}
		j := i + 1
		for j < len(format) && format[j] != '}' {
			j++
		}
		if j >= len(format) {
			b.WriteByte(format[i])
			i++
			continue
		}

		content := format[i+1 : j]
		var idx int
		if len(content) == 0 {
			idx = pos
		} else {
			n := 0
			ok := true
			for k := 0; k < len(content); k++ {
				ch := content[k]
				if ch < '0' || ch > '9' {
					ok = false
					break
				}
				n = n*10 + int(ch-'0')
			}
			if !ok {
				b.WriteString(format[i : j+1])
				i = j + 1
				pos++
				continue
			}
			idx = n
		}

		b.WriteString(get(idx))
		pos++
		i = j + 1
	}
	return b.String()
}
