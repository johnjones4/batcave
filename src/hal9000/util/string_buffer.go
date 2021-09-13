package util

type StringBuffer struct {
	Elements []string
	Pointer  int
}

func NewStringBuffer(size int) StringBuffer {
	return StringBuffer{
		Elements: make([]string, size),
		Pointer:  0,
	}
}

func (b *StringBuffer) Push(e string) {
	b.Elements[b.Pointer%len(b.Elements)] = e
	b.Pointer++
}

func (b StringBuffer) Contains(e string) bool {
	for _, e1 := range b.Elements {
		if e1 == e {
			return true
		}
	}
	return false
}
