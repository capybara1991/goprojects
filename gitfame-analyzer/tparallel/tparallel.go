package tparallel

import "sync"

type T struct {
	parent   *T
	mu       sync.Mutex
	parallel bool
	barrier  chan struct{}
	done     chan struct{}
	wg       sync.WaitGroup
}

func (t *T) Parallel() {
	t.mu.Lock()
	t.parallel = true
	close(t.done)
	t.mu.Unlock()
	<-t.barrier
}

func (t *T) Run(subtest func(t *T)) {
	sub := &T{
		parent:  t,
		barrier: make(chan struct{}),
		done:    make(chan struct{}),
	}

	t.wg.Add(1)
	sub.wg.Add(1)

	go func() {
		subtest(sub)
		sub.wg.Done()

		sub.mu.Lock()
		isParallel := sub.parallel
		sub.mu.Unlock()

		if !isParallel {
			close(sub.done)
		}
	}()

	<-sub.done

	sub.mu.Lock()
	isParallel := sub.parallel
	sub.mu.Unlock()

	if isParallel {
		go func() {
			defer t.wg.Done()
			<-t.barrier
			close(sub.barrier)
			sub.wg.Wait()
		}()
	} else {
		close(sub.barrier)
		sub.wg.Wait()
		t.wg.Done()
	}
}

func Run(topTests []func(t *T)) {
	var parallelTs []*T

	for _, fn := range topTests {
		t := &T{
			barrier: make(chan struct{}),
			done:    make(chan struct{}),
		}

		t.wg.Add(1)
		go func(f func(*T), tt *T) {
			f(tt)
			tt.wg.Done()

			tt.mu.Lock()
			isParallel := tt.parallel
			tt.mu.Unlock()

			if !isParallel {
				close(tt.done)
			}
		}(fn, t)

		<-t.done

		t.mu.Lock()
		isParallel := t.parallel
		t.mu.Unlock()

		if isParallel {
			parallelTs = append(parallelTs, t)
		} else {
			close(t.barrier)
			t.wg.Wait()
		}
	}

	var wg sync.WaitGroup
	for _, t := range parallelTs {
		wg.Add(1)
		go func(tt *T) {
			defer wg.Done()
			close(tt.barrier)
			tt.wg.Wait()
		}(t)
	}
	wg.Wait()
}
