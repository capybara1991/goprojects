//go:build !solution

package waitgroup

type WaitGroup struct {
	lock  chan struct{}
	done  chan struct{}
	count int
}

func New() *WaitGroup {
	wg := &WaitGroup{
		lock: make(chan struct{}, 1),
		done: make(chan struct{}),
	}
	wg.lock <- struct{}{}
	close(wg.done)
	return wg
}

func (wg *WaitGroup) Add(delta int) {
	if wg.lock == nil {
		wg.lock = make(chan struct{}, 1)
		wg.lock <- struct{}{}
	}
	if wg.done == nil {
		wg.done = make(chan struct{})
		close(wg.done)
	}
	<-wg.lock
	if wg.count == 0 && delta > 0 {
		wg.done = make(chan struct{})
	}
	wg.count += delta
	if wg.count < 0 {
		wg.lock <- struct{}{}
		panic("negative WaitGroup counter")
	}
	if wg.count == 0 {
		close(wg.done)
	}
	wg.lock <- struct{}{}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	if wg.lock == nil {
		wg.lock = make(chan struct{}, 1)
		wg.lock <- struct{}{}
	}
	if wg.done == nil {
		wg.done = make(chan struct{})
		close(wg.done)
	}
	<-wg.lock
	ch := wg.done
	wg.lock <- struct{}{}
	<-ch
}
