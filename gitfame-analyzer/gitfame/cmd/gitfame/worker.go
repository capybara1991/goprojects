package main

import (
	"context"
	"sync"
)

type Result struct {
	File   string
	Blocks []BlameBlock
	Err    error
}

func RunPool(ctx context.Context, repo, sha string, files []string, useCommitter bool, n int) <-chan Result {
	if n <= 0 {
		n = 1
	}
	jobs := make(chan string)
	out := make(chan Result)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for f := range jobs {
			select {
			case <-ctx.Done():
				out <- Result{File: f, Err: ctx.Err()}
				continue
			default:
			}
			blocks, err := BlamePorcelain(repo, sha, f, useCommitter)
			out <- Result{File: f, Blocks: blocks, Err: err}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}

	go func() {
		for _, f := range files {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- f:
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
