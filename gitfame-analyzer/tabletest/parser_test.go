package tabletest

import (
	"strings"
	"testing"
	"time"
)

type want struct {
	d   time.Duration
	err string
}

func TestParseDuration_Table(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want want
	}{
		{name: "zero", in: "0", want: want{d: 0}},
		{name: "minus_zero", in: "-0", want: want{d: 0}},

		{name: "plus_simple", in: "+1s", want: want{d: time.Second}},
		{name: "minus_simple", in: "-2s", want: want{d: -2 * time.Second}},

		{name: "ns", in: "150ns", want: want{d: 150 * time.Nanosecond}},
		{name: "us", in: "3us", want: want{d: 3 * time.Microsecond}},
		{name: "micro_u00b5", in: "4µs", want: want{d: 4 * time.Microsecond}},
		{name: "micro_greek_mu", in: "5μs", want: want{d: 5 * time.Microsecond}},
		{name: "ms", in: "12ms", want: want{d: 12 * time.Millisecond}},
		{name: "s", in: "7s", want: want{d: 7 * time.Second}},
		{name: "m", in: "8m", want: want{d: 8 * time.Minute}},
		{name: "h", in: "9h", want: want{d: 9 * time.Hour}},

		{name: "fraction_seconds", in: "1.5s", want: want{d: time.Second + 500*time.Millisecond}},
		{name: "leading_dot_fraction", in: ".25s", want: want{d: 250 * time.Millisecond}},
		{name: "no_fraction_digits_but_dot_ok", in: "1.s", want: want{d: time.Second}},

		{name: "multi_parts", in: "2h45m", want: want{d: 2*time.Hour + 45*time.Minute}},
		{name: "multi_parts_with_fraction", in: "1m0.5s", want: want{d: time.Minute + 500*time.Millisecond}},

		{name: "empty", in: "", want: want{err: "invalid duration"}},
		{name: "just_sign", in: "+", want: want{err: "invalid duration"}},
		{name: "starts_with_unit", in: "ms", want: want{err: "invalid duration"}},
		{name: "dot_unit_no_digits", in: ".s", want: want{err: "invalid duration"}},
		{name: "missing_unit_after_number", in: "1.5", want: want{err: "missing unit"}},
		{name: "unknown_unit", in: "10q", want: want{err: "unknown unit"}},

		{name: "leading_int_overflow", in: "9223372036854775808ns", want: want{err: "invalid duration"}},
		{name: "pre_mul_overflow", in: "2562048h", want: want{err: "invalid duration"}},
		{name: "sum_overflow", in: "2562047h1h", want: want{err: "invalid duration"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.in)
			if tt.want.err != "" {
				if err == nil {
					t.Fatalf("ParseDuration(%q) error expected %q, got nil", tt.in, tt.want.err)
				}
				if !strings.Contains(err.Error(), tt.want.err) {
					t.Fatalf("ParseDuration(%q) error mismatch: want contains %q, got %q", tt.in, tt.want.err, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDuration(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want.d {
				t.Fatalf("ParseDuration(%q)=%v, want %v", tt.in, got, tt.want.d)
			}
		})
	}
}
func TestParseDuration_Extras(t *testing.T) {
	t.Run("very long fraction equals stdlib", func(t *testing.T) {
		frac := strings.Repeat("9", 500)
		in := "0." + frac + "s"
		got, err := ParseDuration(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := time.ParseDuration(in)
		if err != nil {
			t.Fatalf("stdlib parse error: %v", err)
		}
		if got != want {
			t.Fatalf("got %v, want %v for %q", got, want, in)
		}
	})

	t.Run("negative leading dot fraction", func(t *testing.T) {
		got, err := ParseDuration("-.5h")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != -30*time.Minute {
			t.Fatalf("got %v, want %v", got, -30*time.Minute)
		}
	})

	t.Run("plus zero special case", func(t *testing.T) {
		got, err := ParseDuration("+0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0 {
			t.Fatalf("got %v, want 0", got)
		}
	})

	t.Run("invalid starting space", func(t *testing.T) {
		if _, err := ParseDuration(" 1s"); err == nil {
			t.Fatalf("expected error for input with leading space")
		}
	})
}
