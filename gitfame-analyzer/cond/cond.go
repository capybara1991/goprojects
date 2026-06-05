//go:build !solution

package cond

type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L   Locker
	qmu chan struct{}
	q   []chan struct{}
}

func New(l Locker) *Cond {
	c := &Cond{L: l, qmu: make(chan struct{}, 1)}
	c.qmu <- struct{}{}
	return c
}

func (c *Cond) Wait() {
	ch := make(chan struct{})
	<-c.qmu
	c.q = append(c.q, ch)
	c.qmu <- struct{}{}
	c.L.Unlock()
	<-ch
	c.L.Lock()
}

func (c *Cond) Signal() {
	<-c.qmu
	if len(c.q) > 0 {
		ch := c.q[0]
		c.q = c.q[1:]
		close(ch)
	}
	c.qmu <- struct{}{}
}

func (c *Cond) Broadcast() {
	<-c.qmu
	for _, ch := range c.q {
		close(ch)
	}
	c.q = nil
	c.qmu <- struct{}{}
}
