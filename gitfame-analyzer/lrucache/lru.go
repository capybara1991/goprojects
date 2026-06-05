//go:build !solution

package lrucache

import "container/list"

type entry struct {
	key, value int
}

type lru struct {
	cap int
	ll  *list.List
	m   map[int]*list.Element
}

func New(cap int) Cache {
	if cap < 0 {
		cap = 0
	}
	return &lru{
		cap: cap,
		ll:  list.New(),
		m:   make(map[int]*list.Element, cap),
	}
}

func (c *lru) Get(key int) (int, bool) {
	if el, ok := c.m[key]; ok {
		c.ll.MoveToBack(el)
		return el.Value.(*entry).value, true
	}
	return 0, false
}

func (c *lru) Set(key, value int) {
	if el, ok := c.m[key]; ok {
		el.Value.(*entry).value = value
		c.ll.MoveToBack(el)
		return
	}
	if c.cap == 0 {
		return
	}
	el := c.ll.PushBack(&entry{key: key, value: value})
	c.m[key] = el
	if c.ll.Len() > c.cap {
		front := c.ll.Front()
		if front != nil {
			en := front.Value.(*entry)
			delete(c.m, en.key)
			c.ll.Remove(front)
		}
	}
}

func (c *lru) Range(f func(key, value int) bool) {
	for e := c.ll.Front(); e != nil; e = e.Next() {
		en := e.Value.(*entry)
		if !f(en.key, en.value) {
			return
		}
	}
}

func (c *lru) Clear() {
	c.ll.Init()
	for k := range c.m {
		delete(c.m, k)
	}
}
