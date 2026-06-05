//go:build !solution

package otp

import (
	"io"
)

type streamReader struct {
	r    io.Reader
	prng io.Reader
	ks   []byte
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &streamReader{r: r, prng: prng}
}

func (s *streamReader) Read(p []byte) (int, error) {
	n, err := s.r.Read(p)
	if n <= 0 {
		return n, err
	}
	if cap(s.ks) < n {
		s.ks = make([]byte, n)
	} else {
		s.ks = s.ks[:n]
	}
	if err2 := readFullNoErr(s.prng, s.ks); err2 != nil {
		return 0, err2
	}

	for i := 0; i < n; i++ {
		p[i] ^= s.ks[i]
	}
	return n, err
}

type streamWriter struct {
	w    io.Writer
	prng io.Reader
	buf  []byte
	ks   []byte
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &streamWriter{w: w, prng: prng}
}

func (s *streamWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if cap(s.buf) < len(p) {
		s.buf = make([]byte, len(p))
	} else {
		s.buf = s.buf[:len(p)]
	}
	if cap(s.ks) < len(p) {
		s.ks = make([]byte, len(p))
	} else {
		s.ks = s.ks[:len(p)]
	}
	if err := readFullNoErr(s.prng, s.ks); err != nil {
		return 0, err
	}
	for i := range p {
		s.buf[i] = p[i] ^ s.ks[i]
	}
	total := 0
	for total < len(s.buf) {
		n, err := s.w.Write(s.buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	return len(p), nil
}
func readFullNoErr(r io.Reader, b []byte) error {
	off := 0
	for off < len(b) {
		n, err := r.Read(b[off:])
		off += n
		if err != nil {
			return err
		}
	}
	return nil
}
