package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// a sequence of array values are specified with:
// a dash, whitespace, the value, trailing whitespace.
// then loops back to itself to handle the next dash.
type Sequence struct {
	indent int
	parent CollectionTarget
	values []any // tbd: possibly a pointer to the slice?
}

func NewSequence(tgt CollectionTarget, indent int) *Sequence {
	return &Sequence{parent: tgt, indent: indent}
}

func (p *Sequence) History(indent int) (ret CollectionTarget) {
	if p.indent == indent {
		ret = p
	} else {
		ret = p.parent
	}
	return
}

func (p *Sequence) WriteValue(val any) (_ error) {
	p.values = append(p.values, val)
	return
}

func (p *Sequence) NewRune(r rune) (ret charm.State) {
	if r != SequenceDash {
		p.parent.WriteValue(p.values)
	} else {
		lede := Whitespace{Indent: p.indent, required: true}
		ret = charm.Step(&lede, charm.Statement("seq lede", func(r rune) (ret charm.State) {
			// after a dash, the end of file means we're done
			if r == charm.Eof {
				p.WriteValue(nil)
				ret = p.NewRune(r) // writes values
			} else if lede.Indent <= p.indent {
				// after a dash, the same or a lesser amount of whitespace than our own means
				// we must have changed lines: we're either at a new sequence value, or some new parent value.
				p.WriteValue(nil)
				ret = PopHistory(p, lede.Indent).NewRune(r)
			} else {
				// an increased amount of indentation means a value.
				// ( including possibly a sub-sequence )
				ret = charm.RunStep(r, NewValue(p, lede.Indent), charm.Statement("seq value", func(r rune) (ret charm.State) {
					var tail Whitespace
					return charm.RunStep(r, &tail, charm.Statement("seq tail", func(r rune) (ret charm.State) {
						if tail.Lines == 0 && r != charm.Eof { // there should be a newline after the value.
							e := errutil.New("invalid character after sequence value")
							ret = charm.Error(e)
						} else {
							// loop! we're at the next char after a value;
							// by default, we assume the next dash of this same sequence
							// ( a sub-sequence would have been handled as the value )
							ret = p.NewRune(r)
						}
						return
					}))
				}))
				return
			}
			return
		}))
	}
	return
}
