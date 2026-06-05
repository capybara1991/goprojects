//go:build !solution

package rwmutex

type RWMutex struct {
	mu      chan struct{}
	room    chan struct{}
	readers int
}

func New() *RWMutex {
	rw := &RWMutex{
		mu:   make(chan struct{}, 1),
		room: make(chan struct{}, 1),
	}
	rw.mu <- struct{}{}
	rw.room <- struct{}{}
	return rw
}

func (rw *RWMutex) RLock() {
	<-rw.mu
	if rw.readers == 0 {
		<-rw.room
	}
	rw.readers++
	rw.mu <- struct{}{}
}

func (rw *RWMutex) RUnlock() {
	<-rw.mu
	rw.readers--
	if rw.readers == 0 {
		rw.room <- struct{}{}
	}
	rw.mu <- struct{}{}
}

func (rw *RWMutex) Lock() {
	<-rw.room
}

func (rw *RWMutex) Unlock() {
	rw.room <- struct{}{}
}
