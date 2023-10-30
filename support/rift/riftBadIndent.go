package rift

type badIndent struct {
	depth int // number of spaces
}

func (badIndent) Error() string {
	return "bad indent"
}
