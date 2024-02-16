package grok

type Word struct {
	hash  uint64
	slice string // go doesn't allocate a new string for a slice, it stores offset and length
	start int    // but doesn't seem to be a way of finding the offset from the string itself
}

func (w *Word) isValid() bool {
	return len(w.slice) > 0
}

func (w *Word) equals(other uint64) bool {
	return w.hash == other
}

func (w *Word) Equals(other uint64) bool {
	return w.hash == other
}

func (w *Word) Hash() uint64 {
	return w.hash
}

func (w *Word) String() string {
	return w.slice
}
