//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"errors"
	"io"
	"os"
	"sort"
)

type lr struct {
	r   *bufio.Reader
	eof bool
}

func NewReader(r io.Reader) LineReader {
	return &lr{r: bufio.NewReader(r)}
}

func (x *lr) ReadLine() (string, error) {
	if x.eof {
		return "", io.EOF
	}
	s, err := x.r.ReadString('\n')
	if errors.Is(err, io.EOF) {
		x.eof = true
		if len(s) == 0 {
			return "", io.EOF
		}
	} else if err != nil {
		return "", err
	}
	if n := len(s); n > 0 && s[n-1] == '\n' {
		s = s[:n-1]
	}
	return s, nil
}

type lw struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) LineWriter {
	return &lw{w: bufio.NewWriter(w)}
}

func (x *lw) Write(l string) error {
	if _, err := x.w.WriteString(l); err != nil {
		return err
	}
	if _, err := x.w.WriteString("\n"); err != nil {
		return err
	}
	return x.w.Flush()
}

type heapItem struct {
	val string
	i   int
}

type minH []heapItem

func (h minH) Len() int            { return len(h) }
func (h minH) Less(i, j int) bool  { return h[i].val < h[j].val }
func (h minH) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minH) Push(x interface{}) { *h = append(*h, x.(heapItem)) }
func (h *minH) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &minH{}
	heap.Init(h)

	for i, r := range readers {
		if r == nil {
			continue
		}
		line, err := r.ReadLine()
		if err == nil {
			heap.Push(h, heapItem{val: line, i: i})
		} else if !errors.Is(err, io.EOF) {
			return err
		}
	}

	for h.Len() > 0 {
		it := heap.Pop(h).(heapItem)
		if err := w.Write(it.val); err != nil {
			return err
		}
		if line, err := readers[it.i].ReadLine(); err == nil {
			heap.Push(h, heapItem{val: line, i: it.i})
		} else if !errors.Is(err, io.EOF) {
			return err
		}
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	var filesToClose []*os.File
	defer func() {
		for _, f := range filesToClose {
			_ = f.Close()
		}
	}()

	readers := make([]LineReader, 0, len(in))

	for _, path := range in {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		filesToClose = append(filesToClose, f)
		r := NewReader(f)

		var lines []string
		for {
			s, err := r.ReadLine()
			if err == nil {
				lines = append(lines, s)
				continue
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		sort.Strings(lines)

		if err := f.Close(); err != nil {
			return err
		}
		fw, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err = fw.Seek(0, 0); err != nil {
			_ = fw.Close()
			return err
		}
		wr := NewWriter(fw)
		for _, s := range lines {
			if err := wr.Write(s); err != nil {
				_ = fw.Close()
				return err
			}
		}
		if err := fw.Close(); err != nil {
			return err
		}
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		filesToClose = append(filesToClose, fr)
		readers = append(readers, NewReader(fr))
	}

	return Merge(NewWriter(w), readers...)
}
