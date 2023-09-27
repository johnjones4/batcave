package util

import (
	"errors"
	"io"
)

type BytesWriteSeeker struct {
	buf []byte
	pos int
}

func NewBytesWriteSeeker() *BytesWriteSeeker {
	return &BytesWriteSeeker{
		buf: make([]byte, 0),
		pos: 0,
	}
}

func (m *BytesWriteSeeker) Bytes() []byte {
	b := make([]byte, len(m.buf))
	copy(b, m.buf)
	return b
}

func (m *BytesWriteSeeker) Write(p []byte) (n int, err error) {
	minCap := m.pos + len(p)
	if minCap > cap(m.buf) { // Make sure buf has enough capacity:
		buf2 := make([]byte, len(m.buf), minCap+len(p)) // add some extra
		copy(buf2, m.buf)
		m.buf = buf2
	}
	if minCap > len(m.buf) {
		m.buf = m.buf[:minCap]
	}
	copy(m.buf[m.pos:], p)
	m.pos += len(p)
	return len(p), nil
}

func (m *BytesWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = m.pos + offs
	case io.SeekEnd:
		newPos = len(m.buf) + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	m.pos = newPos
	return int64(newPos), nil
}
