//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

type Limiter struct {
	req      chan request
	cancel   chan chan error
	stop     chan struct{}
	done     chan struct{}
	interval time.Duration
	max      int
}

type request struct {
	ctx  context.Context
	resp chan error
}

var ErrStopped = errors.New("limiter stopped")

func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		req:      make(chan request),
		cancel:   make(chan chan error),
		stop:     make(chan struct{}, 1),
		done:     make(chan struct{}),
		interval: interval,
		max:      maxCount,
	}
	go l.loop()
	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.done:
		return ErrStopped
	default:
	}
	resp := make(chan error, 1)
	r := request{ctx: ctx, resp: resp}
	select {
	case l.req <- r:
	case <-l.done:
		return ErrStopped
	}
	go func() {
		select {
		case <-ctx.Done():
			select {
			case l.cancel <- resp:
			case <-l.done:
			}
		case <-resp:
		case <-l.done:
		}
	}()
	select {
	case err := <-resp:
		return err
	case <-l.done:
		return ErrStopped
	}
}

func (l *Limiter) Stop() {
	select {
	case l.stop <- struct{}{}:
	default:
	}
	<-l.done
}

func (l *Limiter) loop() {
	defer close(l.done)
	var (
		queue  []request
		times  []time.Time
		timer  *time.Timer
		timerC <-chan time.Time
	)

	resetTimer := func(d time.Duration) {
		if d <= 0 {
			d = time.Nanosecond
		}
		if timer == nil {
			timer = time.NewTimer(d)
			timerC = timer.C
			return
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(d)
	}

	stopTimer := func() {
		if timer != nil {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer = nil
			timerC = nil
		}
	}

	cleanup := func(now time.Time) {
		cut := now.Add(-l.interval)
		i := 0
		for i < len(times) && times[i].Before(cut) {
			i++
		}
		if i > 0 {
			times = times[i:]
		}
	}

	tryGrant := func(now time.Time) {
		cleanup(now)
		for len(queue) > 0 && len(times) < l.max {
			r := queue[0]
			queue = queue[1:]
			select {
			case r.resp <- nil:
			default:
			}
			times = append(times, now)
			now = time.Now()
			cleanup(now)
		}
		if len(queue) > 0 && len(times) >= l.max {
			wait := l.interval - now.Sub(times[0])
			resetTimer(wait)
		} else if len(queue) == 0 {
			stopTimer()
		}
	}

	for {
		now := time.Now()
		cleanup(now)
		select {
		case r := <-l.req:
			if len(queue) == 0 && len(times) < l.max {
				select {
				case r.resp <- nil:
				default:
				}
				times = append(times, time.Now())
			} else {
				queue = append(queue, r)
				if len(times) < l.max {
					tryGrant(time.Now())
				} else {
					wait := l.interval - now.Sub(times[0])
					resetTimer(wait)
				}
			}
		case ch := <-l.cancel:
			i := -1
			for idx, rr := range queue {
				if rr.resp == ch {
					i = idx
					break
				}
			}
			if i >= 0 {
				rr := queue[i]
				copy(queue[i:], queue[i+1:])
				queue = queue[:len(queue)-1]
				select {
				case rr.resp <- rr.ctx.Err():
				default:
				}
				if len(queue) == 0 {
					stopTimer()
				}
			}
		case <-timerC:
			tryGrant(time.Now())
		case <-l.stop:
			for _, r := range queue {
				select {
				case r.resp <- ErrStopped:
				default:
				}
			}
			return
		}
	}
}
