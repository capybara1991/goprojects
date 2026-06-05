//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type op func(*Evaluator) error

type Evaluator struct {
	stack []int
	dict  map[string]op
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{dict: make(map[string]op)}
	e.dict["+"] = func(e *Evaluator) error {
		a, b, ok := e.pop2()
		if !ok {
			return errors.New("arity")
		}
		e.stack = append(e.stack, a+b)
		return nil
	}
	e.dict["-"] = func(e *Evaluator) error {
		a, b, ok := e.pop2()
		if !ok {
			return errors.New("arity")
		}
		e.stack = append(e.stack, a-b)
		return nil
	}
	e.dict["*"] = func(e *Evaluator) error {
		a, b, ok := e.pop2()
		if !ok {
			return errors.New("arity")
		}
		e.stack = append(e.stack, a*b)
		return nil
	}
	e.dict["/"] = func(e *Evaluator) error {
		a, b, ok := e.pop2()
		if !ok {
			return errors.New("arity")
		}
		if b == 0 {
			return errors.New("div0")
		}
		e.stack = append(e.stack, a/b)
		return nil
	}
	e.dict["dup"] = func(e *Evaluator) error {
		if len(e.stack) < 1 {
			return errors.New("arity")
		}
		e.stack = append(e.stack, e.stack[len(e.stack)-1])
		return nil
	}
	e.dict["drop"] = func(e *Evaluator) error {
		if len(e.stack) < 1 {
			return errors.New("arity")
		}
		e.stack = e.stack[:len(e.stack)-1]
		return nil
	}
	e.dict["swap"] = func(e *Evaluator) error {
		if len(e.stack) < 2 {
			return errors.New("arity")
		}
		n := len(e.stack)
		e.stack[n-1], e.stack[n-2] = e.stack[n-2], e.stack[n-1]
		return nil
	}
	e.dict["over"] = func(e *Evaluator) error {
		if len(e.stack) < 2 {
			return errors.New("arity")
		}
		e.stack = append(e.stack, e.stack[len(e.stack)-2])
		return nil
	}
	return e
}

func (e *Evaluator) Process(row string) ([]int, error) {
	toks := strings.Fields(row)
	for i := 0; i < len(toks); i++ {
		t := toks[i]
		if t == ":" {
			if i+1 >= len(toks) {
				return nil, errors.New("bad def")
			}
			name := strings.ToLower(toks[i+1])
			if _, ok := parseInt(name); ok {
				return nil, errors.New("number redefine")
			}
			j := i + 2
			def := []string{}
			for j < len(toks) && toks[j] != ";" {
				def = append(def, toks[j])
				j++
			}
			if j >= len(toks) || toks[j] != ";" {
				return nil, errors.New("no ;")
			}
			compiled, err := e.compile(def)
			if err != nil {
				return nil, err
			}
			e.dict[name] = func(ev *Evaluator) error {
				for _, f := range compiled {
					if err := f(ev); err != nil {
						return err
					}
				}
				return nil
			}
			i = j
			continue
		}
		if v, ok := parseInt(t); ok {
			e.stack = append(e.stack, v)
			continue
		}
		fn, ok := e.dict[strings.ToLower(t)]
		if !ok {
			return nil, errors.New("unknown")
		}
		if err := fn(e); err != nil {
			return nil, err
		}
	}
	cp := make([]int, len(e.stack))
	copy(cp, e.stack)
	return cp, nil
}

func (e *Evaluator) compile(def []string) ([]op, error) {
	out := make([]op, 0, len(def))
	for _, tok := range def {
		if v, ok := parseInt(tok); ok {
			val := v
			out = append(out, func(ev *Evaluator) error {
				ev.stack = append(ev.stack, val)
				return nil
			})
			continue
		}
		fn, ok := e.dict[strings.ToLower(tok)]
		if !ok {
			return nil, errors.New("unknown")
		}
		out = append(out, fn)
	}
	return out, nil
}

func parseInt(s string) (int, bool) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return int(n), true
}

func (e *Evaluator) pop2() (int, int, bool) {
	if len(e.stack) < 2 {
		return 0, 0, false
	}
	n := len(e.stack)
	a := e.stack[n-2]
	b := e.stack[n-1]
	e.stack = e.stack[:n-2]
	return a, b, true
}
