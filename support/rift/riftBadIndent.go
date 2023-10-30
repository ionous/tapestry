package rift

type badIndent struct {
	have, want int // number of spaces
}

func (badIndent) Error() string {
	return "bad indent"
}
