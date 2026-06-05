//go:build !solution

package speller

import (
	"strings"
)

var ones = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
var tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
var thousands = []string{"", "thousand", "million", "billion"}

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return "minus " + Spell(-n)
	}
	parts := []string{}
	for i := 0; n > 0 && i < len(thousands); i++ {
		chunk := int(n % 1000)
		if chunk != 0 {
			chunkWords := spellChunk(chunk)
			if thousands[i] != "" {
				chunkWords = chunkWords + " " + thousands[i]
			}
			parts = append([]string{chunkWords}, parts...)
		}
		n /= 1000
	}
	return strings.Join(parts, " ")
}

func spellChunk(n int) string {
	res := []string{}
	if n >= 100 {
		res = append(res, ones[n/100]+" hundred")
		n %= 100
		if n == 0 {
			return strings.Join(res, " ")
		}
	}
	if n >= 20 {
		res = append(res, tens[n/10])
		if n%10 != 0 {
			res[len(res)-1] += "-" + ones[n%10]
		}
	} else if n > 0 {
		res = append(res, ones[n])
	}
	return strings.Join(res, " ")
}
