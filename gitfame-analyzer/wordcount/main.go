//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	maper := make(map[string]int)
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Usage: wordcount [word]\n")
	}
	args := os.Args[1:]
	for _, word := range args {
		file, err := os.ReadFile(word)
		check(err)
		lines := strings.Split(string(file), "\n")
		for _, line := range lines {
			maper[line]++
		}
	}
	for line, c := range maper {
		if c >= 2 {
			fmt.Printf("%d\t%s\n", c, line)
		}
	}

}
