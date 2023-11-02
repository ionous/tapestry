package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
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

func (n *Sequence) append(val any) {
	n.values = append(n.values, val)
}

func (n *Sequence) rewrite(val any) {
	n.values[len(n.values)-1] = val
}

func (n *Sequence) NewRune(first rune) (ret charm.State) {
	if first == SequenceDash {
		// cheating a bit here:
		// if next is only whitespace or an eof
		// there's no hook to write; so add it here.
		// alt: push back a pending state into the history, or track a bool
		n.append(nil)
		startingIndent := n.h.CurrentIndent() // padding is the space between the dash and any value
		ret = RequireSpaces("padding", startingIndent, func(padding int) (ret charm.State) {
			switch {
			// if the indent is less or equal,
			// than we're a new line and the value was null
			case padding <= startingIndent:
				ret = n.h.PopIndent(padding)
			// an increased indentation means a value:
			default:
				// first the value:
				// note: the value can be a sub-sequence
				ret = charm.Step(NewValue(n.h, padding, func(v any) (_ error) {
					n.rewrite(v)
					return
				}), charm.Statement("after value", func(r rune) charm.State {
					// after value we require a newline:
					return charm.RunState(r, RequireLines("tail", padding, func(tail int) (ret charm.State) {
						switch {
						case tail <= startingIndent:
							ret = n.h.PopIndent(tail)
						default:
							// the only valid increases are collections,
							// and after a collection, we should be at a less or equal indent.
							e := errutil.New("invalid indentation")
							ret = charm.Error(e)
						}
						return
					}))
				}))
			}
			return
		})
	}
	return
}
