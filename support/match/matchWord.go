package match

type Word struct {
	hash  uint64
	slice string
}

func MakeWord(slice string) Word {
	return Word{Hash(slice), slice}
}

func (w *Word) Hash() uint64 {
	return w.hash
}

func (w *Word) String() string {
	return w.slice
}
