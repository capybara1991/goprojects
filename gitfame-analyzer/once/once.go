//go:build !solution

package once

type Once struct {
	once chan struct{}
	done chan struct{}
}

func New() *Once {
	o := &Once{
		once: make(chan struct{}, 1),
		done: make(chan struct{}),
	}
	o.once <- struct{}{}
	return o
}

func (o *Once) Do(f func()) {
	select {
	case <-o.once:
		defer func() {
			close(o.done)
			if r := recover(); r != nil {
				panic(r)
			}
		}()
		f()
	default:
		<-o.done
	}
}
