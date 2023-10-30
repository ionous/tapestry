package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// a sequence of array values are specified with:
// a dash, whitespace, the value, trailing whitespace.
// then loops back to itself to handle the next dash.
//
// like yaml, whitespace can be any combination of ascii space or ascii tab:
// https://yaml.org/spec/1.2.2/#white-space-characters
type Sequence struct {
	indent int
	values []Value
}

func NewSequence(indent int) *Sequence {
	return &Sequence{indent: indent}
}

func (p *Sequence) StateName() string {
	return "sequence"
}

func (p *Sequence) GetSequence() (ret []Value, err error) {
	ret = p.values
	return
}

// implements valueState for sub-sequences
func (p *Sequence) GetValue() (ret any, err error) {
	ret = p.values
	return
}

func (p *Sequence) NewRune(r rune) (ret charm.State) {
	if r == SequenceDash {
		lede := Whitespace{Indent: p.indent}
		ret = charm.Step(&lede, charm.Statement("seq lede", func(r rune) (ret charm.State) {
			// because some amount of whitespace is required
			// if the indent is unchanged, we know we're on a new line,
			// and all that's allowed is another sequence entry.
			if lede.Indent == p.indent {
				//  `-...\n-` is okay, it indicates a blank value.
				// `-...\n5` is not okay,
				p.values = append(p.values, Value{})
				ret = p.NewRune(r)
			} else if lede.Indent < p.indent {
				// FIX: de-indent to previous sequence
				e := badIndent{lede.Indent, p.indent}
				ret = charm.Error(e)
			} else {
				// some amount of indentation means a value.
				// ( including possibly a sub-sequence )
				value := ValueParser{indent: lede.Indent}
				ret = charm.RunStep(r, &value, charm.Statement("seq value", func(r rune) (ret charm.State) {
					if val, e := value.GetValue(); e != nil {
						ret = charm.Error(e)
					} else {
						tail := Whitespace{Indent: lede.Indent}
						ret = charm.RunStep(r, &tail, charm.Statement("seq tail", func(r rune) (ret charm.State) {
							_, lines := tail.GetSpacing()
							if lines == 0 && r != charm.Eof {
								e := errutil.New("invalid character after sequence value")
								ret = charm.Error(e)
							} else /* if lines > 0 && depth != p.indent {
								// FIX: de-indent to previous sequence
								e := badIndent{depth, p.indent}
								ret = charm.Error(e)
							} else */{
								// loop! we're at the next non-whitespace char after a value;
								// for example, the next dash of *this* sequence.
								// ( note: a sub-sequence would have been handled *as* the value )
								p.values = append(p.values, val)
								ret = p.NewRune(r)
							}
							return
						}))
					}
					return
				}))
				return
			}
			return
		}))
	}
	return
}
