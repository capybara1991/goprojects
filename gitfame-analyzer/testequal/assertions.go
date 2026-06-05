//go:build !solution

package testequal

import (
	"bytes"
	"fmt"
)

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if eq(expected, actual) {
		return true
	}
	t.Errorf("not equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %s",
		expected, actual, formatMsg(msgAndArgs...))
	return false
}

func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !eq(expected, actual) {
		return true
	}
	t.Errorf("equal:\n\texpected: %v\n\tactual  : %v\n\tmessage : %s",
		expected, actual, formatMsg(msgAndArgs...))
	return false
}

func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !AssertEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !AssertNotEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func formatMsg(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	return fmt.Sprintf(fmt.Sprint(msgAndArgs[0]), msgAndArgs[1:]...)
}

func eq(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch x := a.(type) {
	case int:
		return eqSame[int](x, b)
	case int8:
		return eqSame[int8](x, b)
	case int16:
		return eqSame[int16](x, b)
	case int32:
		return eqSame[int32](x, b)
	case int64:
		return eqSame[int64](x, b)

	case uint:
		return eqSame[uint](x, b)
	case uint8:
		return eqSame[uint8](x, b)
	case uint16:
		return eqSame[uint16](x, b)
	case uint32:
		return eqSame[uint32](x, b)
	case uint64:
		return eqSame[uint64](x, b)
	case uintptr:
		return eqSame[uintptr](x, b)

	case string:
		return eqSame[string](x, b)

	case []byte:
		return eqBytes(x, b)
	case []int:
		return eqIntSlice(x, b)

	case map[string]string:
		return eqStringMap(x, b)
	}

	return false
}

func eqSame[T comparable](x T, b interface{}) bool {
	y, ok := b.(T)
	return ok && x == y
}

func eqBytes(x []byte, b interface{}) bool {
	y, ok := b.([]byte)
	if !ok {
		return false
	}
	if (x == nil) != (y == nil) {
		return false
	}
	if x == nil && y == nil {
		return true
	}
	if len(x) == 0 && len(y) == 0 {
		return false
	}
	return bytes.Equal(x, y)
}

func eqIntSlice(x []int, b interface{}) bool {
	y, ok := b.([]int)
	if !ok {
		return false
	}
	if (x == nil) != (y == nil) {
		return false
	}
	if x == nil && y == nil {
		return true
	}
	if len(x) == 0 && len(y) == 0 {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func eqStringMap(x map[string]string, b interface{}) bool {
	y, ok := b.(map[string]string)
	if !ok {
		return false
	}
	if (x == nil) != (y == nil) {
		return false
	}
	if x == nil && y == nil {
		return true
	}
	if len(x) == 0 && len(y) == 0 {
		return false
	}
	if len(x) != len(y) {
		return false
	}
	for k, vx := range x {
		if vy, ok := y[k]; !ok || vy != vx {
			return false
		}
	}
	return true
}
