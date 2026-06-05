package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

type Progress struct {
	total   int64
	done    int64
	enabled bool
	start   time.Time
}

func NewProgress(total int, enabled bool) *Progress {
	return &Progress{
		total:   int64(total),
		enabled: enabled,
		start:   time.Now(),
	}
}

func (p *Progress) Step(file string) {
	if !p.enabled {
		return
	}
	d := atomic.AddInt64(&p.done, 1)
	fmt.Fprintf(os.Stderr, "[%d/%d] %s\n", d, p.total, file)
}

func (p *Progress) Done() {
	if !p.enabled {
		return
	}
	el := time.Since(p.start).Round(time.Millisecond)
	fmt.Fprintf(os.Stderr, "done in %s\n", el)
}
