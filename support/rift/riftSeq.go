package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

// a sequence of array values are specified with:
// a dash, whitespace, the value, trailing whitespace.
// then loops back to itself to handle the next dash.
type Sequence struct {
	h      *History
	values []any // tbd: possibly a pointer to the slice?
}

// maybe h is a factory even?
func NewSequence(h *History, indent int, writeBack func(vs []any) error) charm.State {
	seq := &Sequence{h: h}
	return h.PushIndent(indent, seq, func() error {
		return writeBack(seq.values)
	})
}

func (n *Sequence) NewRune(first rune) (ret charm.State) {
	if first == SequenceDash {
		// cheating a bit here: if the next rune is whitespace or an eof
		// then there will be no pop() to write a value, so write it here
		// alt: push back a pending state into the history, or track a bool
		n.values = append(n.values, nil)
		ret = parseCollection(n.h, func(val any) (_ error) {
			// and rewrite the value here:
			n.values[len(n.values)-1] = val
			return
		})
	}
	return
}
